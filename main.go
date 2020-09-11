package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/davecgh/go-spew/spew"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var isInsideSnap = os.Getenv("SNAP") != ""

func validAPIToken(token string) bool {
	return validation.Validate(token, is.PrintableASCII) == nil && len(token) > 8
}

func browseTo(url string) error {
	var cmd string
	var args []string

	if isInsideSnap {
		return errors.New("inside sandboxed environment")
	} else {
		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default: // "linux", "freebsd", "openbsd", "netbsd"
			cmd = "xdg-open"
		}
		args = append(args, url)
		return exec.Command(cmd, args...).Start()
	}
}

var dbpath string

func withDBPath(f func(path string)) error {
	w := errWrapper("error grabbing database path")
	if isInsideSnap {
		dbpath = os.Getenv("SNAP_USER_DATA") + "/"
	} else if user, err := user.Current(); err != nil {
		return w(err, "couldn't grab the running user")
	} else if runtime.GOOS == "windows" {
		dbpath = user.HomeDir + `\AppData\Roaming\unofficialkyc\`
	} else {
		dbpath = user.HomeDir + "/.local/share/unofficialkyc/"
	}
	f(dbpath + "local.db")
	return nil
}

//I wish golang had better typing :/
func withDBPathErr(f func(path string) error) error {
	var ranErr error
	if err := withDBPath(func(path string) {
		ranErr = f(path)
	}); err != nil {
		return err
	}
	return ranErr
}

var db *gorm.DB

func withDB(f func(*gorm.DB)) error {
	if db == nil {
		w := errWrapper("error initializing local db")
		if err := withDBPathErr(func(path string) error {
			var err error
			if db, err = gorm.Open(sqlite.Open(path), &gorm.Config{}); err != nil {
				return w(err, "error opening local db")
			} else if err = db.AutoMigrate(&User{}); err != nil {
				return w(err, "error migrating user table for local db")
			} else if err = db.AutoMigrate(&Config{}); err != nil {
				return w(err, "error migrating config table for local db")
			}
			return nil
		}); err != nil {
			return w(err)
		}
	}
	f(db)
	return nil
}

func withDBErr(f func(*gorm.DB) error) error {
	var ret error
	if err := withDB(func(db *gorm.DB) {
		ret = f(db)
	}); err != nil {
		return err
	}
	return ret
}

type User struct {
	gorm.Model
	Name     string
	ApiToken string `gorm:"column:api_token"`
}

func (u *User) PostForm(uri string, vals url.Values) (*http.Response, error) {
	var ret *http.Response
	err := withConfigErr(
		func(conf *Config) error {
			if req, err := http.NewRequest("POST", conf.ApiEndpoint+uri, bytes.NewBuffer([]byte(vals.Encode()))); err != nil {
				return err
			} else {
				req.Header.Set("Authorization", u.ApiToken)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				ret, err = http.DefaultClient.Do(req)
				return err
			}
		},
	)
	return ret, err
}

func withUser(f func(user *User)) error {
	return withConfigErr(func(conf *Config) error {
		w := errWrapper("error grabbing logged in user from db")
		if conf.User.Name == "" {
			fmt.Println("Haven't authenticated yet; please log in.")
			fmt.Print("Username: ")
			var username string
			fmt.Scanln(&username)
			fmt.Print("Password: ")
			if password, err := secureTermRead(); err != nil {
				return w(err, "error reading password from terminal")
			} else if resp, err := http.PostForm(conf.ApiEndpoint+"/new_api_token", url.Values{
				"username": []string{username},
				"password": []string{password},
			}); err != nil {
				return w(err, "error requesting refresh token from api")
			} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
				return w(err, "error reading api login response body")
			} else if resp.StatusCode == http.StatusUnauthorized {
				return errors.New("Incorrect username or password")
			} else if resp.StatusCode != 200 {
				return w(errors.New(string(b)), "api returned non-200 response code when trying to get a new API token, along with the following body")
			} else if token := string(b); !validAPIToken(token) {
				return w(errors.New(string(b)), "the api returned a success status code, but the following, structurally invalid api token")
			} else {
				conf.User.Name = username
				conf.User.ApiToken = strings.TrimSpace(string(b))
				if err := db.Save(&conf.User).Error; err != nil {
					return w(err, "error saving user into configuration")
				} else {
					conf.UserID = conf.User.ID
					if err := db.Save(&conf).Error; err != nil {
						return w(err, "error saving user into configuration")
					}
				}
			}
		}
		f(&conf.User)
		return nil
	})
}

type Config struct {
	gorm.Model
	ApiEndpoint string `gorm:"column:api_endpoint"`
	UserID      uint
	User        User
}

var conf *Config

func withConfig(f func(*Config)) error {
	if conf == nil {
		w := errWrapper("error getting config from db")
		if err := withDBErr(
			func(db *gorm.DB) error {
				var configs []Config
				if err := db.Preload("User").Find(&configs).Error; err != nil {
					return err
				} else if len(configs) == 0 {
					conf = &Config{
						ApiEndpoint: "https://unofficialkyc.com/api/v1",
					}
					return db.Save(conf).Error
				} else if len(configs) == 1 {
					conf = &configs[0]
					return nil
				} else {
					return errors.New("You have multiple configs in your db. You should probably delete it. It's in local/share/unofficialkyc")
				}
			},
		); err != nil {
			return w(err)
		}
	}
	f(conf)
	return nil
}

func withConfigErr(f func(*Config) error) error {
	var ret error
	if err := withConfig(func(conf *Config) {
		ret = f(conf)
	}); err != nil {
		return err
	}
	return ret
}

func secureTermRead() (string, error) {
	if b, err := terminal.ReadPassword(syscall.Stdin); err != nil {
		return "", err
	} else {
		fmt.Println()
		return string(b), err
	}
}

func printHelp() {
	fmt.Println(`
    List of commands:
    register - Registers a new UFKYC passport.
    token - Grab a UFKYC token for the domain in your clipboard.
    donate [amount] - Donate to add to your credibility score (and buy some Kenyan kid a malaria net).
    service register - Registers a UFKYC service users will be able to generate.
    service register_domain [name] - Adds an unvalidated domain to your UFKYC service, and starts the validation process.
    service require_donation [amount] - (Optional) Adds an amount users have to have donated in order to create tokens for your service.
    `)
}

func dangerous(f func()) {
	if os.Getenv("DANGEROUS") != "TRUE" {
		fmt.Println("You don't have the DANGEROUS=TRUE environment variable set. This command requires it; please don't use api_switch unless you are either a UFKYC developer or want to get owned.")
	}
	f()
}

//Only braindead monkeys whine about programs putting a lot of code in main().
//Contrary to popular belief, taking your laundry list and dividing it into
//doThis() and doThat() subroutines does not automatically make your code
//cleaner.  If you have a refactoring suggestion make sure it's not that _real_
//dumb one.

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must specify a command.")
		printHelp()
	} else {
		command := os.Args[1]
		switch command {
		case "api_switch":
			dangerous(func() {
				if len(os.Args) == 3 {
					if !strings.HasPrefix(os.Args[2], "http") || validation.Validate(os.Args[2], is.URL) != nil {
						fmt.Println("Passed argument is not a valid URL.")
					} else {
						printErr(withConfig(func(config *Config) {
							config.ApiEndpoint = os.Args[2]
							if err := db.Save(config).Error; err != nil {
								fmt.Println("Error saving new API endpoint into database:", err)
							} else {
								fmt.Println("All your commands will now contact", os.Args[2], "for api requests.")
							}
						}))
					}
				}
			})
		case "register":
			//We need the config so we can save the user and also so that we can know what endpoint to contact
			printErr(withConfig(func(conf *Config) {
				if conf.User.Name != "" {
					fmt.Println("You have already logged in as user '" + conf.User.Name + "'.")
				} else {
					fmt.Print("Username: ")
					var username string
					fmt.Scanln(&username)
					password := ""
					for {
						var err error
						fmt.Print("Password: ")
						if password, err = secureTermRead(); err != nil {
							fmt.Println("Couldn't read password from terminal; " + err.Error())
						} else {
							fmt.Print("Confirm password: ")
							if passwordConfirmation, err := secureTermRead(); err != nil {
								fmt.Println("Couldn't read password from terminal; " + err.Error())
							} else if password != passwordConfirmation {
								fmt.Println("Passwords were not the same, try again: ")
							} else {
								break
							}
						}
					}
					if resp, err := http.PostForm(conf.ApiEndpoint+"/register", url.Values{
						"username": []string{username},
						"password": []string{password},
					}); err != nil {
						fmt.Println("error requesting api:", err)
					} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
						fmt.Println("error reading response body:" + err.Error())
					} else if responseStr := string(b); responseStr == "user already exists" {
						fmt.Println("User already exists; try a different username.")
					} else if resp.StatusCode != http.StatusOK {
						fmt.Println("received non-200 status code and the following response body:", responseStr)
					} else if token := string(b); !validAPIToken(token) {
						fmt.Println("the api returned a success status code, but the following, structurally invalid api token:", token)
					} else {
						conf.User = User{
							Name:     username,
							ApiToken: token,
						}
						if err := withDBErr(func(db *gorm.DB) error {
							return db.Save(conf).Error
						}); err != nil {
							fmt.Println("Registered the user, but there was a problem saving them to the database: ", err)
						}
					}
				}
			}))
		case "donate":
			var amount float64
			philanthropize := func() {
				fmt.Println("Enter an email address to be associated with the payment, in case of disputes. You may use a tempmail if desired:")
				var email string
				for {
					fmt.Scanln(&email)
					if err := validation.Validate(email, is.Email); err != nil {
						fmt.Println("Email entered was invalid; try again:")
					} else {
						break
					}
				}
				printErr(withUser(func(user *User) {
					if resp, err := user.PostForm("/donate", url.Values{
						"amount":         []string{strconv.FormatFloat(amount, 'f', 5, 64)},
						"payment_vendor": []string{"globee"},
						"email":          []string{email},
					}); err != nil {
						fmt.Println("Error contacting API (no payment was made):", err)
					} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
						fmt.Println("Error reading API response (no payment was made):", err)
					} else if resp.StatusCode != http.StatusOK {
						fmt.Println("API returned with an error (no payment was made) and the following response body:", string(b))
					} else if url := strings.TrimSpace(string(b)); validation.Validate(url, is.URL) != nil {
						fmt.Println("Strange; the API returned a non-url to browse to to continue payment, but delivered an OK status code. Here was the URL:")
						fmt.Println(url)
					} else if isInsideSnap {
						clipboard.WriteAll(url)
						fmt.Println("Please browse to the URL pasted into your clipboard and finish your cryptocurrency payment.")
						fmt.Println("Your donation will be confirmed shortly therafter.")
					} else if err := browseTo(url); err != nil {
						fmt.Println("An error occured opening the payment URL: ", err)
						fmt.Println("Please attempt to go to", url, " in whatever browser you have available manually to finish your payment.")
					} else {
						fmt.Println("Please attempt to finish your cryptocurrency payment in the opened browser tab. Your donation will be confirmed shortly thereafter.")
					}
				}))
			}
			if len(os.Args) < 3 {
				fmt.Println("Enter amount you want to amount, in U.S. dollars: ")
				if n, err := fmt.Scanf("%f\n", &amount); n != 1 || err != nil {
					fmt.Println("Couldn't parse payment amount;", err)
				} else {
					philanthropize()
				}
			} else {
				var err error
				if amount, err = strconv.ParseFloat(strings.TrimRight(os.Args[2], "$"), 64); err != nil {
					fmt.Println("An amount argument was provided, but it wasn't a decimal number. Try again.")
				} else {
					philanthropize()
				}
			}
		case "service":
			if len(os.Args) < 3 {
				fmt.Println("Subcommand to 'service' is required (register, etc.)")
				printHelp()
			} else {
				switch os.Args[2] {
				case "register":
					printErr(withUser(func(user *User) {
						if resp, err := user.PostForm("/register_service", url.Values{}); err != nil {
							fmt.Println("Error encountered while contacting api:", err)
						} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
							fmt.Println("Error reading response body:", err)
						} else if respStr := strings.TrimSpace(string(b)); resp.StatusCode != http.StatusOK {
							fmt.Printf("API returned the status code %d and the following response body: %s\n", resp.StatusCode, respStr)
						} else {
							fmt.Println("Your service registration was sucessful, and your service's granted ID is '" + respStr + "'. Assign it some domain names to allow users to generate tokens for it.")
						}
					}))
				case "require_donation":
					if len(os.Args) != 4 {
						fmt.Println("You used the wrong number of arguments; this command needs 4.")
						printHelp()
					} else {
						if amount, err := strconv.ParseFloat(strings.TrimSuffix(os.Args[3], "$"), 64); err != nil {
							fmt.Println("Error parsing donation amount: " + err.Error())
							printHelp()
						} else {
							printErr(withUser(func(user *User) {
								if resp, err := user.PostForm("/require_donation", url.Values{
									"amount": []string{strconv.FormatFloat(amount, 'f', 2, 64)},
								}); err != nil {
									fmt.Println("Error trying to connect to API:", err)
								} else if resp.StatusCode != 200 {
									if b, err := ioutil.ReadAll(resp.Body); err != nil {
										fmt.Println("API returned the status code " + strconv.Itoa(resp.StatusCode) + ".")
									} else {
										fmt.Println("API returned the status code", resp.StatusCode, "and the following response body: "+string(b))
									}
								} else {
									fmt.Printf("New users will now have to donate at least %0.2f$ platform wide in order to start creating tokens for your service.\n", amount)
								}
							}))
						}
					}
				case "register_domain":
					printErr(withUser(func(user *User) {
						do := func(domain string) {
							if resp, err := user.PostForm("/register_service_domain", url.Values{
								"domain_name": []string{domain},
							}); err != nil {
								fmt.Println("Error trying to connect to API:", err)
							} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
								fmt.Println("Error trying to read API response body:", err)
							} else if resp.StatusCode != http.StatusOK {
								var errMsg struct {
									Error string `json:"error"`
								}
								if err := json.Unmarshal(b, &errMsg); err != nil {
									fmt.Println("API returned non-200 status code, and we were unable to unmarshal the error message. Here it is raw: " + strings.TrimSpace(string(b)))
									fmt.Println("And here's the error encountered during unmarshaling:", err)
								} else {
									fmt.Println("The API returned an error: " + errMsg.Error)
								}
							} else {
								var resp struct {
									Data struct {
										PathValidation struct {
											Path    string `json:"path"`
											Content string `json:"content"`
										} `json:"path_validation"`
										TxtValidation struct {
											Nonce string `json:"nonce"`
										} `json:"txt_validation"`
									} `json:"data"`
								}
								fmt.Println(string(b))
								if err := json.Unmarshal(b, &resp); err != nil {
									fmt.Println("The API returned with a success, but we were unable to marshal the response. Here is what it sent us, raw: " + spew.Sdump(resp))
								} else {
									if resp.Data.PathValidation.Content != "" {
										fmt.Println("Your domain name has been registered.")
										fmt.Println("In order to validate ownership, you'll need to place a file at the '" + resp.Data.PathValidation.Path + "' path of a web server running on port 80 or 443.")
										fmt.Println("The file must contain the following nonce: '" + resp.Data.PathValidation.Content + "'")
										fmt.Println("UFKYC will continually poll that location from the internet until it responds correctly, at which point your domain will be validated.")
										fmt.Println("If you do not validate ownership within an hour, your domain will become unregistered and you'll need to start this process again.")
										fmt.Println("You can re-run this command to get the above information again from UFKYC.")
									} else if resp.Data.TxtValidation.Nonce != "" {
										fmt.Println("Your domain name has been registered.")
										fmt.Println("In order to validate ownership, you'll need make a TXT record at the root domain")
										fmt.Println("with the contents '" + resp.Data.TxtValidation.Nonce + "'.")
										fmt.Println("We will continually poll its TXT records until it responds correctly.")
										fmt.Println("If you do not validate ownership within an hour, your domain will become unregistered and you'll need to start this process again.")
										fmt.Println("You can re-run this command to get the above information again from UFKYC.")
									}
								}

							}
						}
						if len(os.Args) == 4 {
							if validation.Validate(os.Args[3], is.Domain) != nil || !isRootDomain(os.Args[3]) {
								fmt.Println("Passed argument is not a valid root domain.")
							} else {
								do(os.Args[3])
							}
						} else {
							var domain string
							for {
								fmt.Print("Enter domain: ")
								fmt.Scanln(&domain)
								if validation.Validate(domain, is.Domain) != nil || !isRootDomain(domain) {
									fmt.Println("Entry was not a valid root domain; try again.")
								} else {
									var confirm string
									fmt.Print("Confirm: ")
									fmt.Scanln(&confirm)
									if confirm != domain {
										fmt.Println("Domain and confirmation were different; try again.")
									} else {
										break
									}
								}
							}
							do(domain)
						}
					}))
				default:
					fmt.Println("Subcommand unrecognized.")
					printHelp()
				}
			}
		case "clear":
			dangerous(func() {
				printErr(withDBPathErr(func(path string) error {
					if err := os.Remove(path); err != os.ErrNotExist {
						return err
					} else {
						return nil
					}
				}))
			})
		case "token":
			printErr(withUser(func(user *User) {
				if clipboard.Unsupported {
					fmt.Println("Sorry, clipboard functionality was not found for your current running environment.")
					if runtime.GOOS == "linux" {
						fmt.Println("Make sure you have the clipboard program installed for your preferred display manager (xclip, xsel, wl-clip, etc.)")
					}
				} else if domain, err := clipboard.ReadAll(); err != nil {
					fmt.Println("We encountered an error reading your clipboard:", err)
				} else if domain = strings.TrimSpace(domain); validation.Validate(domain, is.Domain) != nil || !isRootDomain(domain) {
					fmt.Println("The item in your clipboard was not a domain. Make sure you copy the root domain in your browser before trying to generate a token.")
					fmt.Println("It's a pain, but this way hopefully you'll never get phished again.")
				} else {
					fmt.Println("Grab token for", domain, "(y/n)?")
					if r, _, _ := bufio.NewReader(os.Stdin).ReadRune(); r == 'y' || r == 'Y' {
						if resp, err := user.PostForm("/get_account_token", url.Values{
							"service_domain": []string{domain},
						}); err != nil {
							fmt.Println("Error encountered while contacting api for new token:", err)
						} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
							fmt.Println("Error encountered while reading response body of api request:", err)
						} else if rstr := strings.TrimSpace(string(b)); resp.StatusCode != 200 {
							fmt.Println("The API rejected your request for a token and responded with the following:", rstr+".")
						} else if err := clipboard.WriteAll(rstr); err != nil {
							fmt.Println("Error encountered writing token to clipboard:", err)
						} else {
							fmt.Println("Token copied to clipboard.")
						}
					}
				}
			}))
		default:
			fmt.Println("Command not recognized.")
			printHelp()
		}
	}
}
