recurse() {
    for i in "$1"/*; do
        if [ -d "$i" ]; then
            recurse "$i"
        elif [ -f "$i" ]; then
            if [[ $i = *".sol"* ]]; then
                solc --allow-paths . --overwrite --abi $i -o ./smartcontract/build/abi
                solc --allow-paths . --overwrite --bin $i -o ./smartcontract/build/bin
            fi
        fi
    done
}

pwd
recurse ./smartcontract
