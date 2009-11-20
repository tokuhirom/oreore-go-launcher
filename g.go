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

var archmap = map[string] string {
    "386":  "8",
    "amd64":"6",
    "arm":  "5",
};

func run(tmpl string, v ...) int {
    cmd := fmt.Sprintf(tmpl, v);
    return system.System(cmd);
}

func compile(src string) (string, bool) {
    arch, ok := archmap[os.Getenv("GOARCH")];
    if !ok {
        log.Exit("invalid environment variable $GOARCH");
    }
    obj := func() string {
        if src[len(src)-3:len(src)] == ".go" {
            return src[0:len(src)-2] + arch;
        }
        return src + "." + arch;
    }();
    out      := arch + ".out";

    if run("%s -o %s %s", arch+"g", obj,  src) != 0 {
        return "", false;
    }
    if run("%s -o %s %s", arch+"l", out,  obj) != 0 {
        return "", false;
    }
    return out, true;
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

