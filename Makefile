.PHONEY: test_ls

test_ls:
	go build -o egg cmd/egg/main.go
	./egg ls -al

test_git:
	go build -o egg cmd/egg/main.go
	./egg git status

test_bad_git:
	go build -o egg cmd/egg/main.go
	./egg git asdf
