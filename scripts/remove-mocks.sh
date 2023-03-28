if [[ ! -v TEST_NAME ]]; then
    echo "TEST_NAME is not set"
    exit 1
fi

for file in $(grep -l $TEST_NAME ./mock/mappings/*); do
    echo "Removing $file"
    rm -i -f $file;
    #  ^ prompt for delete
done
