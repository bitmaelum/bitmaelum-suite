[![Go Report Card](https://goreportcard.com/badge/github.com/bitmaelum/bitmaelum-suite)](https://goreportcard.com/report/github.com/bitmaelum/bitmaelum-suite)
[![Build Status](https://travis-ci.org/bitmaelum/bitmaelum-suite.svg?branch=master)](https://travis-ci.org/bitmaelum/bitmaelum-suite)
![License](https://img.shields.io/github/license/bitmaelum/bitmaelum-suite)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bitmaelum/bitmaelum-suite)         
         
         ____  _ _   __  __            _                 
        |  _ \(_) | |  \/  |          | |                
        | |_) |_| |_| \  / | __ _  ___| |_   _ _ __ ___  
        |  _ <| | __| |\/| |/ _` |/ _ \ | | | | '_ ` _ \ 
        | |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |
        |____/|_|\__|_|  |_|\__,_|\___|_|\__,_|_| |_| |_|
           P r i v a c y   i s   y o u r s   a g a i n                                          

# What if..
we could redesign email without any need for backward compatibility. What would it look like? 
Could we solve the problems we face nowadays like spam, email forgery, phishing emails and, maybe the most important 
one: privacy?

BitMaelum (old Anglo-Saxon for "bit-by-bit") is an attempt. Instead of trying to figure out how to fix email, we are building a new email system from the ground up. This way, we can design the system to solve mail problems at its core.

To read the current status of the project, check out [our status document](documentation/status.md).

# Benefits
We try to design a system that:

  - Deal with less to no spam at all
  - Nullifies mail address harvesting (it becomes useless)
  - Host your mail wherever you like
  - Move your messages to another provider without losing your mail address
  - Allows that only you can subscribe and unsubscribe from mailing lists
  - Slow connection / mobile friendly 
  - Privacy first

# How do we want to achieve all this?
There are a lot of issues we want to solve. Some of them more difficult than others. In a nutshell, this is what BitMaelum does:
 
### Deal with less to no spam at all
There are multiple ways we try to combat spam:

  - We verify that the email came from the original sender and has not been tampered with in any way.
  - Sending out emails costs effort, making it less economical for spammers to send out thousands (even millions) of messages. When botnets are used (computers that have been taken over that sends out spam), it will trigger enough workload to become either noticeable for the botnet victim.
  - When sending legitimate emails, we allow users to enroll in mailing-lists. This gives the mailing-list owner a personal key that will enable you to send mail to users for this specific list. Even if this key gets public through hacks and leaks, no other organization or user cannot do anything with this information, as this key is not valid for anyone else.
  - Not being signed up for this list is impossible? Since signing up becomes an action from the recipient and not from the sender, this makes the recipient in charge. Don't want to receive any more emails? Just revoke your key for that mailing list, and they are not able to send you mail. 
  - Selling email addresses is not viable anymore. Personal keys are only valid for that specific organization and not for others, making collections of many email addresses pointless.
 

### Nullifies mail address harvesting
We cannot prevent that your email address will be leaked. In fact, assume it will, no matter how careful you are.

But when your address is out there in the open, it will not be feasible for spam lists to send you mail. Legitimate mailing lists can only mail you with a personal key that only works for that specific mailing list. So even you can send somebody mail, sending thousands of millions of people mail is practically not possible. 


### Host your mail wherever you like
Setting up a mail server is easy and by default, secure out of the box. Just point your mail address to your new server, and you are good to go! 
You can host your mail at home, at your office, or as most users do, at a mail provider.


### move your messages to another mailserver without loosing your mail address
Don't like your provider anymore, or are they getting too pricy? Or maybe you changed your internet provider. BitMaelum allows you to take your message AND your address with you. You are not bound to any email provider anymore!

 
### Allows that only you can subscribe and unsubscribe from mailing lists
Tired of getting email you never subscribed to? Or ever tried to unsubscribe to such a list, only to find yourself getting MORE mail?

BitMaelum puts you in charge of mailing lists. You and ONLY you decide if you want to subscribe to a mailing list. Want to unsubscribe? This is also YOUR decision. There is no need to wait for any action from the mailing list owner. As soon as you decide to unsubscribe, you're unsubscribed!
 
### Slow connection / mobile friendly
Not everybody has a 10Gbit fiber internet connection. And in the middle of the woods, you're probably lucky with a 1-bar connection.
We try to keep communication to a minimum by selective downloading what you need.

### Privacy first
We use end-to-end encryption between the sender and recipient. When sending a message, ONLY the recipient can decrypt the email. This means that no server, router, or agency that is snooping traffic can read your messages or email meta-data. Even your mail provider cannot read or modify your messages.

The only piece of information that is known publicly is the hashed of the sender and recipient (so you can't even detect to who we send or receive a message from). 


# Are you guys serious!?
We don't expect this to hit the production level ever. It's mostly a nice hobby-project to identify issues with current email infrastructure and see if we can come up with (realistic) solutions. Will they work? Do the benefits outweigh using a completely new system? It can.. but we don't assume it will. Still, it's a fun project.


![https://bitmaelum.com/logo_and_name.svg](https://bitmaelum.com/logo_and_name.svg)
