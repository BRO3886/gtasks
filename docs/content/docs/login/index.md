---
title: 'Logging In'
draft: false
weight: 3
summary: Login to Google Tasks CLI with your Google Account
---

* To login, all you need to do is type:
```
gtasks login
```
* This will open show you an output like this:
```
* Go to the following link in your browser then type the authorization code: 
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=something.apps.googleusercontent.com&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ftasks&state=state-token

Enter the code: 
```
* If you open the link in your browser, it will ask you to sign in with your Google account. **If it shows you that the app is not verified, you can click Advanced Options and continue with the login process** 
* Once you accept all the permissions it is asking (such as creating and managing tasks), Google will generated a code for you. Come back to the terminal and paste the code. Your login should be complete and you'll be able to see and manage your tasks.
