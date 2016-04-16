# mailrandom
Mail yourself a random selection from a list of weighted items

I schedule this program to run daily at 3 a.m., so that a randomly selected activity is waiting in my inbox when I wake up.

## Usage examples

To specify weighted items using the command line and have a randomly selected item mailed to you, run:

    mailrandom -o "practice spanish","play saxophone","read philosophy" -w 1,1,2 -m mailconfig-example.json -p my_mail_password

The weights specified by the `-w` flag make "read philosophy" twice as likely to be chosen as either "practice spanish" or "play saxophone." The file `mailconfig-example.json` indicates the format of the mail configuration file. If SMTP authentication is not required, the `Username` field may be omitted from the mail configuration file.

To randomly choose an item from a CSV file containing a list of weighted items, run:

    mailrandom -i activities-example.csv -m mailconfig-example.json

If SMTP authentication is required but a password is not specified using the `-p` flag, the program will look for the password in the `MAIL_PASS` environment variable.

To see all command line options, run:

    mailrandom -h

## Weighting

If item k has weight `w_k`, it will be selected with probability `w_k / (w_1 + w_2 + ... + w_total)`. Weights can be non-negative integers or floating point numbers.
