package main
import (
    "os";
    "log";
    "system";
    "readline";
    "fmt";
    "io";
    "strings";
);

func compile(src string) (string, bool) {
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
    obj := func() string {
        if src[len(src)-3:len(src)] == ".go" {
            return src[0:len(src)-2] + arch;
        }
        return src + ".go";
    }();
    linker := arch + "l";
    compiler := arch + "g";
    out := arch + ".out";

    if system.System(compiler + " -o " + obj + " " +  src + " && " + linker + " -o " + out + " " + obj) == 0 {
        return out, true;
    }
    return "", false;
}

func run_interpreter() {
    out, succeeded := compile(os.Args[1]);
    if succeeded {
        os.Exec(out, os.Args[2:len(os.Args)], os.Environ());
    }
}

func run_shell() {
    prompt := "go> ";

    for {
        result := readline.ReadLine(&prompt);
        switch {
        case result == nil:
            return;
        case *result == "":
            // nop. ignore empty line
        default:
            run_one_line(*result);
            readline.AddHistory(*result);
        }
    }
}

func run_one_line(line string) {
    src := render_src(line);
    fname := tempfile();
    defer os.Remove(fname);
    io.WriteFile(fname, strings.Bytes(src), 0600);
    out, succeeded := compile(fname);
    if succeeded {
        files := make([]*os.File, 3);
        files[0] = os.Stdin;
        files[1] = os.Stdout;
        files[2] = os.Stderr;
        pid, err := os.ForkExec(out, make([]string, 0), os.Environ(), "", files);
        if err != nil {
            log.Stderr(err, "\n");
            return;
        } else {
            os.Wait(pid, 0);
        }
    } else {
        log.Stderr("compilation error\n");
    }
}

func render_src(line string) string {
    src := fmt.Sprintf(`
package main
func main() {
    %s
}
`, line);
    return src;
}

func tempfile() string {
    return fmt.Sprintf("%d.go", os.Getpid());
}

func main() {
    if len(os.Args) < 2 {
        run_shell();
    } else {
        run_interpreter();
    }
}

