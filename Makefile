.PHONEY: test_ls

.DEFAULT_GOAL := clean_egg

egg:
	go build -o egg cmd/egg/main.go

clean:
	rm -f egg

clean_egg: clean egg

test_ls: egg
	./egg ls -al

test_git: egg
	./egg git status

test_bad_git: egg
	./egg git asdf
