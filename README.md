# snailmail

snailmail manages sending and receiving zip files through email.

For most people, if they want to share a file with someone, they can upload it to a cloud storage provider and share the link. But, what about the people that don't have access to cloud storage? They have to resort to email, that's what they have to do. Unfortunately, gmail limits attachment file sizes to 25mb. This is where snailmail comes into play. Given a file, snailmail will break it into 25mb chuncks, zip it up, and then send them all to a recipient. Then, the recipient can use snailmail to download all of the files, join, and unzip them. 

This project is still a work in progress. To test it, create file and run:

```go get```

```go build && snailmail.exe -a=send -e="dave8051gmail.com" -f="test.txt" ```