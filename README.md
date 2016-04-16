# mailrandom
Mail yourself a random selection from a list of weighted items

I schedule this program to run daily at 3 a.m., so that a randomly selected exercise is waiting in my inbox when I wake up.

## Usage examples

Weighted items specified using command line:
    mailrandom -o "practice spanish","play saxophone","read philosophy" -w 1,1,2 -m mailconfig-example.json -p my_mail_password
The weights specified by the `-w` flag make "read philosophy" twice as likely to be chosen as either "practice spanish" or "play saxophone."

Weighted items specified in a CSV file:
    mailrandom -i activities-example.csv -m mailconfig-example.json
Since an SMTP mail server password is not specified, the program will look to the `MAIL_PASS` environment variable. If SMTP authentication is not required, omit the `Username` field from the mail configuration JSON file.

## Weighting

If item k has weight `w_k`, it will be selected with probability `w_k / (w_1 + w_2 + ... + w_total)`. Weights can be non-negative integers or floating point numbers.
