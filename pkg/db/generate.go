package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i Transactor -o ./mocks -s "_minimock.go"
