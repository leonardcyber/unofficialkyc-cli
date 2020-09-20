# unofficialKYC Command Line Interface

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/kycli)

This repository houses the GPLv3 licensed UnofficialKYC CLI. People can use
this application to create tokens, accounts, and services on the UnofficialKYC
platform, as well as increase their UFKYC credibility scores by donating money
or cryptocurrency to the third world.

This project is in the very early stages of its existence. We're publishing it
because we want you to use it, but be aware, you're probably going to have
personal conversations with the devs in our Keybase (link at the bottom) about
bugs you find or ambiguous setup steps. We're happy to help you manually
install and work through issues you find using it so that it will begin to get
tested.

## Table of Contents
<!-- vim-markdown-toc GFM -->

* [FAQ](#faq)
    * [What is unofficialKYC?](#what-is-unofficialkyc)
    * [Who am I paying and why? You say it's going to charity?](#who-am-i-paying-and-why-you-say-its-going-to-charity)
    * [How does UFKYC's authentication flow work?](#how-does-ufkycs-authentication-flow-work)
    * [Why did you start by publishing a CLI and not a GUI or web application?](#why-did-you-start-by-publishing-a-cli-and-not-a-gui-or-web-application)
    * [What payment methods does UFKYC accept?](#what-payment-methods-does-ufkyc-accept)
* [Installation](#installation)
    * [Snap](#snap)
    * [Manual compilation (it's not as hard as it usually is)](#manual-compilation-its-not-as-hard-as-it-usually-is)
* [Support](#support)

<!-- vim-markdown-toc -->
## FAQ

### What is unofficialKYC?

unofficialKYC is a platform for users to link their accounts on many online
services to a single charitable donation, without having to share their
disclose their real life identity or shared identity between those services. In
doing this, UFKYC provides a way for companies to penalize the creation of many
different accounts, without requiring genuinely new users to pay money they
haven't already provided for another service, and without having to require
identifying information like passports or IDs.

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
requiring users to create a "stake", a given amount donated on behalf of a
UFKYC passport to charity along with a small service fee. Users who have
donated at or above the amount specified by the service they're using are
allowed to generate "authentication tokens", which, depending on the service,
give access to the site or are used as apart of the authentication process.

The beauty of this system is that you can pay the money once, and then generate
a token for any service that accepts signups at that stake level or higher.
This means that bans for traditionally porous internet communities can finally
have "teeth" - regular UFKYC users can sign up for most services with no
paywall beyond their first, but rulebreakers, shills, and trolls are forced to
make new UFKYC passports and provide the stake again and again.  At the same
time, these services will never receive your passport ID, and the tokens you
generate will only tell the service who you are on their platform. Unlike the
OAUTH "sign in with google" buttons, sites you use UFKYC to authenticate to
won't be able to look up what other accounts you have signed in for with the
same account by comparing email addresses.

Because the funds are an arbitrary "proof-of-stake", and exist to enable the
service as much to financially support UFKYC, all funds staked in the way are
donated to [GiveDirectly](https://www.givedirectly.org). GiveDirectly is a
top-rated nonprofit foundation that provides electronic cash transfers to the
world's poorest. GiveDirectly has been a [GiveWell](https://www.givewell.org)
top rated charity for 8 years running, and you should consider donating to them
anyways if you don't use our services.

### How does UFKYC's authentication flow work?

There are two main use scenarios for UFKYC. The first is as a signup procedure,
where users generate tokens for sites once during the account creation process.
The second is as a logon procedure, where users generate tokens each time they
log into a website as a primary or second factor authentication mechanism. Both
scenarios involve the same steps:

1. A service operator registers their service with `kycli service register`, and is given a
   service ID.
2. The service operator and attaches and validates some domains to their
   service via the `kycli service register_domain` command, which will ensure
   users who generate tokens for their domain will have their service ID as an
   audience.
2. The service operator modifies a login, landing, or signup page to include a
   form for user tokens, along with perhaps a link to
   [unofficialkyc.com](https://unofficialkyc.com) as explanation.
3. A user sees the option to authenticate/signup with a UFKYC token and copies
   the root domain from their browser, asks the CLI for a token with `kycli
   token`, and submits it kycli passed back into the clipboad. This
   [PASETO](https://paseto.io/) token passed to the service includes the
   service ID given during service registration, a new, service-specific "subject" identifier
   for that user,
4. The service receives the token, and verifies its signature, audience, and
   subject. Then it persists the token's subject in a database or associates it
   with an identity on the service so that the same user can't sign up twice.
   Now that the user has been authenticated (and inadvertently proven they
   aren't being phished), the service can proceed with account creation or
   login.

### Why did you start by publishing a CLI and not a GUI or web application?

For a few reasons:
1. It's simply easier for us to start by developing a CLI that interfaces with
   an API than a GUI desktop application that does the same, or a standalone
   web application, because we don't have to worry about graphic design.
2. Leonard Cyber developed this product partly for our personal use
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
   actually by preventing surreptitous access.  Because by definition a web
   application downloads a fresh set of code each time it runs, they can be
   sandboxed, but not really "open sourced" without extraneous FSF-approved
   browser addons. From a security perspective, there's not really a whole lot
   of difference between sending new, obfuscated, unreviewed javascript from
   Protonmail on each run and a strong public promise to keep data secret.
   With a command line interface, the code can be GPLv3 licensed and put in a
   public repository, users can check out specific versions, and at least
   anticipate that someone would have said something if it curled your keys to
   a remote server. It's not a perfect remedy, and introduces some new
   complications, but it's ultimately the better tradeoff given that we have
   sandboxing covered.

A web app on unofficialkyc.com is being developed that will eventually
supplement this program in case you don't agree or don't want to install
anything, and in the long term we're looking at making a desktop GUI, too. We
expect the web app should be finished by the end of September, but make no
promises.

### What payment methods does UFKYC accept?

UFKYC accepts both card (through stripe) and several cryptocurrencies (through
globee) for donations. We're working on a way to get donations sent via a
first-party integration with givedirectly.org; for now, you have to donate
through one of these venues and we pass along the money to them every week.

## Installation

We only officially support Linux. You're welcome to use the windows subsystem
for Linux or compile yourself, which should also 'just work', as the CLI has no
outside dependencies besides a web browser that won't be transparently brought
in during the golang compilation process.

### Snap

Snap may not be as bad as you have been led to believe. At the moment we
publish a snapcraft app named "kycli". To install it, after first installing
[snap](https;//snapcraft.io) (if it's not already available on your
distribution), run the terminal command:

`snap install --candidate kycli`

And that's it.

The snap works on more distributions than we could ever expect to suport
directly, updates automatically, provides a standard environment for the thing
to run in and for us to test, is barely bigger than the default binary, and
allows us to somewhat securely sandbox it from the rest of your operating
system without resorting to docker containers.  It's too much of a quick win
for us to start writing .debs and .emerges and .aurs by hand, when we could be
further developing the platform for the people that use it. 

We do understand, however, that automatic updates for some people are a no-no,
even though in this case we think they're a net positive. Those people who
dislike either of those things should go ahead and try:

### Manual compilation (it's not as hard as it usually is)

To manually compile, first install [go](https://golang.org). Simply clone the
repository and run `go build`. You should be left with a working CLI called
`kycli` in your inside the cloned repository, respectively. You can place that
wherever you like, perhaps in `/usr/bin/`.

## Support

More documentation is coming. For more information on what UFKYC does, how to
use it, or how to build it into your service, we have a [Keybase
team](https://keybase.io/team/unofficialkyc) in which we are almost always
active and can answer your questions.
