.PHONEY: test_ls

test_ls:
	go build -o $$ cmd/wrap/main.go
	./$$ ls -al

test_git:
	go build -o $$ cmd/wrap/main.go
	./$$ git status
