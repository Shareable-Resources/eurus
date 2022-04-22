#!/bin/bash
usage()
{
  echo "Usage:          [ -abi ] abi file
                [ -o ] output file
                [ -a ] alias of the variables
                [ -pkg ] package of output file
                [ -t ] type of the smart contract
        "
  exit 2
}


while getopts ":b:o:a:pkg:t:" arg; do
    case "${arg}" in
        b)
            ABI=${OPTARG}
            ;;
        o)
            OUTPUT=${OPTARG}
            ;;
        a)
            ALIAS=${OPTARG}
            ;;
        p)
            PKG=${OPTARG}
            ;;
        t)
            TYPE=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -z ABI ]
then
    echo "-b <value> is missing, which is used for specified the abi file path for the smart contract."
fi

if [ -z OUTPUT ]
then
    echo "-o <value> is missing, which is used for specified the file path output Go file."
fi

if [ -z PKG ]
then
    echo "-p <value> is missing, which is used for specified the package name of the output Go file."
fi

if [ ! -z ABI ] && [ ! -z OUTPUT ] && [ ! -z ALIAS ] && [ ! -z PKG ] && [ ! -z TYPE ]
then
    abigen --abi=ABI --out=OUTPUT --alias=ALIAS --pkg=PKG --type=TYPE
elif [ ! -z ABI ] && [ ! -z OUTPUT ] && [ ! -z PKG ] && [ ! -z TYPE ]
then
    abigen --abi=ABI --out=OUTPUT --pkg=PKG --type=TYPE
elif [ ! -z ABI ] && [ ! -z OUTPUT ] && [ ! -z PKG ] && [ ! -z ALIAS ]
then
    abigen --abi=ABI --out=OUTPUT --pkg=PKG --alias=ALIAS
elif [ ! -z ABI ] && [ ! -z OUTPUT ] && [ ! -z PKG ]
then
    abigen --abi=ABI --out=OUTPUT --pkg=PKG
fi