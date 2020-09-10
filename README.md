# UnofficialKYC Command Line Interface

This repository houses the GPLv3 licensed UnofficialKYC CLI. Users can use this
application to create tokens, accounts, and services on the UnofficialKYC
platform, as well as increase their UFKYC credibility scores by donating
cryptocurrency to impoverished African children.

This project is in the very early stages of its existence. We're publishing it
because we want you to use it, but be aware, you're probably going to have
personal conversations with the devs in our Keybase (link at the bottom) about
bugs you find or ambiguous setup steps. We're happy to help you manually
install it so it can start to get tested.

## FAQ

### What is UnofficialKYC?

UnofficialKYC is a platform for users to anonymously link their accounts on
many online services to a single charitable donation, without having to share
their identity between those services. In doing this, UFKYC provides a way for
companies to penalize the creation of many different accounts, without
requiring genuinely new users to pay money they haven't already provided for
another service before, and without having to require identifying information
from new users like passports or IDs.

By incorporating UnofficialKYC into their authentication or signup procedures, websites can:

1. Prevent the vast majority of phishing attacks, even for services where the
   user is originally unsure of the domain name they're supposed to
   authenticate to, through our unique authentication workflow.
2. Prohibitively increase the difficulty of launching successful large-scale spam or DDoS attacks.
3. Provide a second factor of authentication that is both automatically backed
   up and does not require per-site configuration by users.
4. Greatly cut down on duplicate account creation, and provide a genuine financial
   hurdle for banned or punished users' re-registeration, without harming your
   good users with a (local and unique) paywall.
5. Help reduce extreme poverty.

### Who am I paying and why? You say it's going to charity?

The services registered with UFKYC prevent redundant account creation by
requiring users to front a "stake", a given amount donated on behalf of a UFKYC
passport to charity. Users who have donated at or above the amount specified by
the service they're using are allowed to generate "authentication tokens",
which, depending on the service, give access to the site or are used as apart
of the authentication process.

The beauty of this system is that you can donate the money once, and then
generate a token for any service that allows users from that stake level or
higher. This means that bans for traditionally porous internet communities can
finally have "teeth" - regular UFKYC users can sign up for most services with
no paywall beyond an initial payment, but rulebreakers, shills, and trolls are
forced to make new UFKYC passports and provide the stake again and again.  At
the same time, these services will never receive your passport ID, and the
tokens you generate will only tell the service who you are on their platform
and won't automatically disclose all of your other accounts.

Because the funds are an arbitrary "proof-of-stake", and exist to enable the
service as much to financially support UFKYC, only 10% of stakes users provide
are actually go to the UFKYC devs. The other 90% are donated to
[GiveDirectly](https://www.givedirectly.org). GiveDirectly is a top-rated
nonprofit foundation that provides electronic cash transfers to the world's
poorest. They are probably the world's most effective and well-researched
charity available today, and you should consider donating to them anyways if
you don't use our services.

### How does UFKYC's authentication flow work?

There are two main use scenarios for UFKYC. The first is as a signup procedure,
where users generate tokens for sites once during the account creation process.
The second is as a logon procedure, where users generate tokens each time they
log into a website as a primary or second factor authentication mechanism. Both
scenarios involve the same core procedure:

1. A service operator registers their service with UFKYC, and adds and validates some domains to the service.
2. The service operator modifies a login, landing, or signup page to include a form for user tokens.
3. A user sees the option to authenticate/signup with a UFKYC token and copies
   the root domain shown, asks the CLI for a token with `kycli token`, and
   pastes the (now copied) token into the field.
4. The service accepts the token, verifies its signature, audience, and
   subject, persists the subject in a database or associates it with an
   identity on the service, and proceeds with account creation or login.

### Why did you start by publishing a CLI and not a GUI or web application?

For a few reasons:
1. It's simply easier to start by developing a CLI that interfaces with an API
   than a GUI desktop application that does the same, or a standalone web
   application, because we don't have to worry about graphic design.
2. We (Leonard Cyber) developed this product partly for our personal use
   administering Leonard's online computer hacking exam, and our users are
   (hopefully) comfortable with CLIs.
3. CLIs are badass.
4. There are certain guarantees we can give with an open source command line
   interface that it is impossible to give via a web application, which we plan
   on pursuing in the near future. One of these is automated public/private key
   management, key publication, and end-to-end encrypted messaging. While
   browsers are perfectly capable of performing encryption, contrary to the
   aspirations of Protonmail or Mega.nz, they can only assure confidentiality
   in a limited way, primarily for the sake of scaring off subpeonas and not
   actually preventing surreptitous access.  Because by definition a web
   application downloads a fresh set of code each time it runs, even if you
   audited the javascript to verify Protonmail wasn't sending the contents of
   your emails back today, you can't be sure it won't the next time you visit
   the site. Whereas, with a command line interface, the code can be GPLv3
   licensed and open-sourced, users can check out specific versions, and users
   can install through software repositories that have people look over the
   code first. It's not a perfect remedy, and introduces some new
   complications, but it's ultimately the better tradeoff from a security
   perspective.

A web app on unofficialkyc.com is being developed that will eventually
supplement this program, and in the long term we're looking at making a desktop
GUI, too. We expect the web app should be finished by the end of September, but
make no promises.

### Why cryptocurrency?

It was easier to setup, we like cryptocurrency, and good ones like Monero allow
for easy, private payments. Existing solutions in this space like per-service
paywalls or traditional KYC usually inadvertently require that users give up
their personal information to a random internet service, and we wanted a
product which didn't require that of our users. That being said, we *are*
working on "regular" payments as an option for users who don't want to have to
buy monopoly money.

## Installation

We only officially support Linux. You're welcome to use the windows subsystem
for Linux or compile yourself, which should also 'just work', as the CLI has no
outside dependencies besides a web browser that won't be transparently brought
in during the golang compilation process.

### Snapcraft

Snap is not as bad as you may have been led to believe. At the moment we
publish a snapcraft app named "kycli". To install it, after first installing
[snap](https;//snapcraft.io) (if it's not already available on your
distribution), run the terminal command:

`snap install kycli`

And that's it.

The snap works on more distributions than we could ever expect to suport
directly, updates automatically, only occupies a trivial amount more space than
it does when the CLI is compiled by hand, provides a standard environment for
the thing to run in and for us to test, and allows us to sandbox it reasonably
securely from the rest of your operating system without resorting to docker
containers. It's too much of a quick win for us to start writing .debs and
.emerges and .aurs by hand at the moment.  Perhaps someday this CLI will be
fantastically popular and this installation guide will be filled with
distro-specific packages, but we're not going to go make them, because there's
the snap.

Nevertheless, we understand automatic updates for some people are a no-no, even
though in this case we think they're a net positive, and those people should go
ahead and try:

### Manual compilation

To manually compile, first install [go](https://golang.org). Then either run
`go install github.com/unofficialkyc-cli`, or simply clone the repository and
run `go build`. You should be left with a working CLI called `kycli` in your
`~/go/bin/` folder or inside the cloned repository, respectively. You can place
that wherever you like.

---

For more information on what UFKYC does, how to use it, or how to build it into
your service, check out [unofficialkyc.com](https://unofficialkyc.com) (If it
is up quite yet). If you want to talk to the devs directly, we have a [Keybase
team](https://keybase.io/team/unofficialkyc).
