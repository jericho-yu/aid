if [ "$1" = "publish" ]; then
  git add --all && \
  git commit -m"$2"

  if [ -n "$3" ]; then
    git push origin $3 && \
  	git checkout master && \
  	git merge $3 && \
  	git push origin master && \
  	git checkout $3
  fi

  if [ -n "$4" ]; then
  	git tag "$4" && git push origin $4
  fi
fi

if [ "$1" == "get-tag" ]; then
  git describe --tags $(git rev-list --tags --max-count=1)
fi