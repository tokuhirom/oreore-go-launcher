package main
import (
    "os";
    "log";
    "system";
);

func main() {
    if len(os.Args) < 2 {
        log.Exit("Usage: ", os.Args[0], " src.go");
    }

    arch := func () string {
        archtype := os.Getenv("GOARCH");
        switch {
        case archtype=="386":
            return "8"
        case archtype=="amd64":
            return "6";
        case archtype=="arm":
            return "5";
        case archtype=="":
            log.Exit("Cannot get a environment variable $GOARCH");
        default:
            log.Exit("unknown archtecture : ", archtype);
        };
        return ""; // should not reach here
    }();
    src := os.Args[1];
    obj := func() string {
        println(src);
        if src[len(src)-3:len(src)] == ".go" {
            return src[0:len(src)-2] + arch;
        }
        return src + ".go";
    }();
    linker := arch + "l";
    compiler := arch + "g";
    out := arch + ".out";

    if system.System(compiler + " " +  src + " && " + linker + " " + obj) == 0 {
        os.Exec(out, os.Args[2:len(os.Args)], os.Environ());
    }
}

