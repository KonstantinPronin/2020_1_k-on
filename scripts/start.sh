mkdir -p build/
now=$(date "+%F-%T")
build_dir="build"
build_filename="go_${now}"
build_path="${build_dir}/${build_filename}"
go build -o $build_path ./cmd/main
cd $build_dir
./$build_filename &