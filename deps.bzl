def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def _kingpin(go_repository):
    _maybe(
        go_repository,
        name = "com_github_alecthomas_kingpin",
        importpath = "github.com/alecthomas/kingpin",
        tag = "v2.2.6",
    )

    _maybe(
        go_repository,
        name = "com_github_alecthomas_units",
        commit = "f65c72e2690dc4b403c8bd637baf4611cd4c069b",
        importpath = "github.com/alecthomas/units",
    )

    _maybe(
        go_repository,
        name = "com_github_alecthomas_template",
        commit = "fb15b899a75114aa79cc930e33c46b577cc664b1",
        importpath = "github.com/alecthomas/template",
    )

def _prometheus(go_repository):
    _maybe(
        go_repository,
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        tag = "v1.6.0",
    )

    _maybe(
        go_repository,
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        tag = "v0.9.1",
    )

    _maybe(
        go_repository,
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        version = "v2.1.1",
        sum = "h1:6MnRN8NT7+YBpUIWxHtefFZOKTAPgGjpQSxqLNn0+qY=",
    )

    _maybe(
        go_repository,
        name = "com_github_beorn7_perks",
        importpath = "github.com/beorn7/perks",
        tag = "v1.0.1",
    )

    _maybe(
        go_repository,
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        commit = "v0.2.0",
    )

    _maybe(
        go_repository,
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        tag = "v0.0.11",
    )

    _maybe(
        go_repository,
        name = "com_github_matttproud_golang_protobuf_extensions",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        commit = "c182affec369e30f25d3eb8cd8a478dee585ae7d",
    )

    _maybe(
        go_repository,
        name = "org_golang_x_sys",
        commit = "1151b9dac4a98d49ef7f80f07ddd826ff51e0b36",
        importpath = "golang.org/x/sys",
    )

def deps(go_repository):
    _maybe(
        go_repository,
        name = "com_github_google_gopacket",
        importpath = "github.com/google/gopacket",
        commit = "a83d5be3a7a49b54dbe6a0c55fd10abd833f3aa2",
    )

    _maybe(
        go_repository,
        name = "com_github_cenkalti_backoff",
        importpath = "github.com/cenkalti/backoff",
        tag = "v4.0.2",
    )

    _kingpin(go_repository)
    _prometheus(go_repository)

def test(go_repository):
    _maybe(
        go_repository,
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        tag = "v0.4.1",
    )
