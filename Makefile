.PHONEY: test_ls


clean:
	rm -f egg

egg:
	go build -o egg cmd/egg/main.go

test_ls: egg
	./egg ls -al

test_git: egg
	./egg git status

test_bad_git: egg
	./egg git asdf
