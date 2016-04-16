# mailrandom
Mail yourself a random selection from a list of weighted items

I schedule this program to run daily at 3 a.m., so that a randomly selected activity is waiting in my inbox when I wake up.

## Usage examples

To specify weighted items using the command line, run:

    mailrandom -o "practice spanish","play saxophone","read philosophy" -w 1,1,2

A randomly selected item will be printed to the console. In this case, the weights specified by the `-w` flag make "read philosophy" twice as likely to be chosen as either "practice spanish" or "play saxophone." 

To have the randomly selected item mailed to you, use the `-m` flag and provide a mail configuration file (see mailconfig-example.json for an example):

    mailrandom -o "jumping jacks","sit ups","push ups" -w 1,1,1 -m mailconfig-example.json -p my_mail_password

The SMTP authentication password can be provided using the `-p` flag, as above, or by setting the `MAIL_PASS` environment variable. If SMTP authentication is not required, omit the `Username` field from the mail configuration file.

To randomly choose an item from a CSV file containing a list of weighted items, run:

    mailrandom -i activities-example.csv -m mailconfig-example.json

To see all command line options, run:

    mailrandom -h

## Weighting

If item k has weight `w_k`, it will be selected with probability `w_k / (w_1 + w_2 + ... + w_total)`. Weights can be non-negative integers or floating point numbers.
