# goReferrer

for [oneplus invitation list](https://oneplus.net/invites)

creates a list of emails based on one email, just seperated by dots then sends out confirmation emails for each but they all go to the same email (it works like this on gmail, any dots in the user part of the email eg hel.lo@gmail.com is equal to hello@gmail.com, it puts a period every 2 characters because when it was one they started blocking it iirc).

second goroutine listens on the inbox for emails and extracts the confirmation link and spoofs it. Then it deletes that email. 

Both also spoof as the browser by copying the exact request that chrome would use for the GET requests except with the path replaced.

When it gets blocked from anything it just increases the timeout between every request by 1, for every failed request. Otherwise for every 10 successful requests it decreases the timeout between requests. That way you can just let it run without having to modify anything yourself, its completely self adjusting. I made this in less than a day, so its not the best code, I see lots of improvements but no point in doing it now that there is the captcha and the gmail trick is blocked for now. 

Well it was tons of fun and it only ran for a few hours and I got a really high spot so I'm happy :) Figured id open source it for anyone else doing this sort of stuff in the future in golang! 

<img src="https://www.imageupload.co.uk/images/2015/08/10/ScreenShot2015-08-10at1.06.26PM.png" border="0">

###SUCCESS - UPDATE
<img src="http://bit.ly/1F4fPO5" border="0">
