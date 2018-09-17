cwd=$(pwd)

build_module () {
    cd $1
    if [ ! -d "lib" ]; then
        mkdir lib && cd lib
    else
        cd lib
    fi
    go build ..
    cd $cwd
}

build_module src/buildings/building_construction