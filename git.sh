operator="$1"

if [ "$operator" = "push-tag" ]; then
	tag="$2"
	git tag "$tag" && git push origin "$tag"
elif [ "$operator" = "last-tag" ]; then
	git describe --tags $(git rev-list --tags --max-count=1)
else
	branch="$2"
	commit="$3"
	tag="$4"

	git add --all &&
		git commit -m"$commit"

	if [ -n "$branch" ]; then
		git push origin $branch &&
			git checkout $operator &&
			git merge $branch &&
			git push origin $operator &&
			git checkout $branch
	fi

	if [ -n "$tag" ]; then
		git tag "$tag" && git push origin "$tag"
	fi
fi