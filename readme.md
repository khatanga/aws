#Khatanga AWS Automation

## Use a credentials file 

To be able to execute against an aws account provide credentials `~/.aws/credentials`

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

To fork and change code to ensure go get still works as expected do this

1. Fork this repo
2. Use go get to clone the original

```
go get github.com/khatanga/aws
```

3. Add your fork as a remote
```
git remote add fork https://github.com/{you}/aws.git
```

4. This added your fork as a remote, which means that once you've modified some files and committed the changes you can now run:

```
git push fork
```

5. Submit pull request