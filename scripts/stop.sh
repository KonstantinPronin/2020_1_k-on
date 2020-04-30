# shellcheck disable=SC2006
id1=`pgrep go_main_service`
id2=`pgrep go_auth_service`
id3=`pgrep go_film_service`
id4=`pgrep go_series_service`

kill $id1
kill $id2
kill $id3
kill $id4
