mkdir -p build/
build_dir="build"

build_filename1="go_main_service"
build_filename2="go_auth_service"
build_filename3="go_film_service"
build_filename4="go_series_service"

build_path1="${build_dir}/${build_filename1}"
build_path2="${build_dir}/${build_filename2}"
build_path3="${build_dir}/${build_filename3}"
build_path4="${build_dir}/${build_filename4}"

go build -o $build_path1 ./cmd
go build -o $build_path2 ./application/microservices/auth/cmd/
go build -o $build_path3 ./application/microservices/film/cmd/
go build -o $build_path4 ./application/microservices/series/cmd/

./$build_path4 &
./$build_path3 &
./$build_path2 &
./$build_path1 &