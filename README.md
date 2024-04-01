# JLPT Practice Helper
### A golang application to help practice for the JLPT exam. 

It builds a daily Japanese lesson and emails it for study.
- it fetches content from tables in Airtable: 
    - kanji
    - vocabulary
    - grammar
    - link to graded reader book on tadoku.org
    - link to youtube video for listening practice
    
- it populates an email template with the information it fetched

- it emails it according to the github actions schedule cron

Environment variables that it needs: 
- EMAIL_FROM
- EMAIL_TO
- EMAIL_CC
- SMTP_USERNAME
- SMTP_PASSWORD
- SMTP_HOST
- SMTP_PORT
- AIRTABLE_API_KEY

The resulting email looks like this:

![email-visual](/static/jlpt_email.gif)