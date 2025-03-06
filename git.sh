operator="$1"

if [ "$operator" = "push-tag" ]; then
	tag="$2"
	git tag "$tag" && git push origin "$tag"
elif [ "$operator" = "last-tag" ]; then
	git describe --tags $(git rev-list --tags --max-count=1)
else
	IFS=':' read -r src_branch dst_branch <<<"$operator"
	commit="$2"
	tag="$3"

	git add --all && git commit -m"$commit"

	if [ -n "$src_branch" ]; then
		git push origin $src_branch
	fi

	if [ -n "$dst_branch" ]; then
		git checkout $dst_branch &&
			git merge $src_branch &&
			git push origin $dst_branch &&
			git checkout $src_branch
	fi

	if [ -n "$tag" ]; then
		git tag "$tag" && git push origin "$tag"
	fi
fi
