# UnofficialKYC Command Line Interface

This repository houses the GPLv3 licensed UnofficialKYC CLI. Users can use this
application to create tokens, accounts, and services on the UnofficialKYC
platform, as well as increase their UFKYC credibility scores by donating
cryptocurrency to impoverished African children.

## FAQ

### What is UnofficialKYC?

UnofficialKYC is a platform for users to link their accounts on many online
services to a single anonymous charitable donation, while maintaining seperate
identities between those services. In doing this, UFKYC provides a way for
companies to penalize the creation of many different accounts, without
requiring genuinely new users to pay money they haven't already provided for
another service before, and without having to require identifying information
from new users like passports or IDs. It also provides a way for those services
to support single sign-on while ensuring the privacy of their users' accounts.

By incorporating UnofficialKYC into their authentication or signup procedures, websites can:

1. Prevent the vast majority of phishing attacks, even for services where the
   user is originally unsure of the domain name they're supposed to authenticate to.
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
passport to charity. Users who have donated the required amount on the platform
are allowed to generate authentication tokens for those services, which,
depending on the service give access to the site or are used as apart of the
authentication process.

The beauty of this system is that you can donate the money once, and then
generate a token for any service that allows users from that stake level or
higher. This means that bans for traditionally porous internet communities can
have "teeth" - regular UFKYC users can sign up for most services with no
paywall beyond the first time, but rulebreakers, shills, and shell users are
forced to make new UFKYC passports and provide the stake again and again. At
the same time, these services don't ever receive your passport ID, and the
tokens you generate will only tell the service who you are on their platform
and won't automatically disclose all of your other accounts.

Because the funds are an arbitrary "proof-of-stake", and exist to enable the
service as much as they pay for UFKYC's development costs, 10% of stakes users
provide are used to pay for UFKYC's development costs and the other 90% provide
are donated to [GiveDirectly](https://www.givedirectly.org). GiveDirectly is a
top-rated nonprofit foundation that provides electronic cash transfers to the
world's poorest. They are probably the world's most effective and
well-researched charity available today, and you should consider donating to
them anyways if you don't use our services.

### How does UFKYC's authentication flow work?

There are two main use scenarios for UFKYC. The first is as a signup procedure, where users
generate tokens for sites once during the account creation process. The second is as a logon
procedure, where users generate tokens each time they log into a website as a primary or second
factor authentication mechanism. Both scenarios involve the same core procedure:

1. A service operator registers their service with UFKYC, and adds and validates some domains to the service.
2. The service operator modifies a login, landing, or signup page to include a form for user tokens.
3. A user sees the option to authenticate/signup with a UFKYC token and copies
   the domain, asks the CLI for one, and pastes it into the field.
4. The service accepts the token and verifies its signature, audience, and
   subject, persists the subject in a database or associates it with an
   identity on the service, and proceeds normally.

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
   can voluntarily choose to wait a few days before applying a new update. It's
   not a perfect remedy, and introduces some new complications, but it's
   ultimately better than the alternative for our use case.

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

---

For more information on what UFKYC does, how to use it, or how to build it into
your service, check out [unofficialkyc.com](https://unofficialkyc.com) (If it
is up quite yet). If you want to talk to the devs directly, we have a [Keybase
team](https://keybase.io/team/unofficialkyc).
