arr=($(cat ./configs/config.yml | grep "password:"))
password=$(eval echo ${arr[1]})

arr=($(cat ./configs/config.yml | grep "path:"))
path=$(eval echo "${arr[1]} ${arr[2]}")

arr=($(cat ./configs/config.yml | grep "dbname:"))
dbname=$(eval echo ${arr[1]})

PGPASSWORD=$password "$path" -h localhost -U postgres -d $dbname \
    -f "./schema/init.sql" \
    -f "./schema/test/admin.sql" \
    -f "./schema/test/user.sql" \
    -f "./schema/test/courier.sql" \
    -f "./schema/test/restaurant.sql" \
    -f "./schema/test/category.sql" \
