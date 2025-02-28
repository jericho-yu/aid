git add --all && \
git commit -m"$1"

if [ -n "$2" ]; then 
  git push origin $2 && \
	git checkout master && \
	git merge $2 && \
	git push origin master && \
	git checkout $2
fi

if [ -n "$3" ]; then
	git tag "$3" && git push origin $3
fi