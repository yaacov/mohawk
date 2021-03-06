%global provider        github
%global provider_tld    com
%global project         yaacov
%global repo            mohawk
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
%global commit          8eefab3f90ab8c828e202a5a0fc20150ecae1ff2
%global shortcommit     %(c=%{commit}; echo ${c:0:7})

Name:           %{repo}
Version:        0.33.3
Release:        1%{?dist}
Summary:        Time series metric data storage
License:        Apache
URL:            https://%{import_path}
Source0:        https://github.com/MohawkTSDB/mohawk/archive/%{version}.tar.gz

BuildRequires:  git
BuildRequires:  golang >= 1.2-7

%description
Mohawk is a metric data storage engine that uses a plugin architecture for data storage and a simple REST API as the primary interface.

%prep
%setup -q -n mohawk-%{version}

# many golang binaries are "vendoring" (bundling) sources, so remove them. Those dependencies need to be packaged independently.
rm -rf vendor

%build
# set up temporary build gopath, and put our directory there
mkdir -p ./_build/src/github.com/MohawkTSDB
ln -s $(pwd) ./_build/src/github.com/MohawkTSDB/mohawk

export GOPATH=$(pwd)/_build:%{gopath}
export GOBIN=$(pwd)/_build/bin
go get ./src
make

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./mohawk %{buildroot}%{_bindir}/mohawk

%files
%defattr(-,root,root,-)
%doc LICENSE README.md
%{_bindir}/mohawk

%changelog
* Tue Jun 26 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.33.3-1
- Do not allow .. in static files path

* Sat Feb 17 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.33.2-1
- Set exports http header to text
- Do not export old data

* Tue Feb 14 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.33.1-1
- Add optional default tenant
- Add Prometheus export endpoint

* Tue Jan 21 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.32.2-1
- Fix duplicate buckets for stats requests

* Tue Jan 11 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.32.1-1
- Add error handlers to api calls
- Alerts accept a list of metrics or a tags regex

* Tue Jan 2 2018 Yaacov Zamir <kobi.zamir@gmail.com> 0.31.1-1
- Storage default to ASC insdead of DESC
- Default limit changed from 2000 to 20000

* Wed Dec 31 2017 Yaacov Zamir <kobi.zamir@gmail.com> 0.30.8-1
- Add options response
- Fix query by tags

* Sat Dec 30 2017 Yaacov Zamir <kobi.zamir@gmail.com> 0.28.2-1
- Response includes feedback
- Remove empty stats response
- Add help (cli) for storage options
- Enable string time format parsing, e.g. -7h, 6mn
- Add the /m metrics endpoint (rest api)
- Add memory storage options, retention and granularity times
- Add mongo storage options, multiple servers, username and password

* Wed Dec 6 2017 Yaacov Zamir <kobi.zamir@gmail.com> 0.27.1-1
- Add min max to memory storage stats

* Wed Dec 6 2017 Yaacov Zamir <kobi.zamir@gmail.com> 0.26.2-7
- Initial RPM release
