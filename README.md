# goReferrer
creates a list of emails based on one email, just seperated by dots then sends out confirmation emails for each but they all go to the same email (it works like this on gmail). second goroutine listens on the inbox for emails and extracts the confirmation link and spoofs it. Then it deletes that email. It also spoofs as a browser by copying the exact request that chrome would use for GET requests. When it gets blocked from anything it just increases the timeout between every request by 1, for every failed request. Otherwise for every 10 successful requests it increases the timeout between requests. That way you can just let it run without having to modify anything yourself. I made this in less than a day, so its not the best code certainly, I see lots of improvements but no point in doing it now. Well it was tons of fun and it only ran for a few hours before they added a captcha lol. Figured id open source it for anyone else doing this sort of stuff in the future in golang! 

<img src="https://www.imageupload.co.uk/images/2015/08/10/ScreenShot2015-08-10at1.06.26PM.png" border="0">
