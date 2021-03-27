extern crate clap;
use clap::{Arg, App, SubCommand};

fn main() {
    let _matches = App::new("yabba")
    .version("0.9.0")
    .author("Deepak Kaul")
    .about("bandwidth analyzer")
    .arg(Arg::with_name("verbose")
        .short("v")
        .long("verbose")
        .help("versbose output")
        .takes_value(false))
    .subcommand(SubCommand::with_name("server")
        .about("runs server")
        .arg(Arg::with_name("port")
            .short("p")
            .long("port")
            .value_name("SERVER_PORT")
            .help("bind port")
            .takes_value(true)))
    .subcommand(SubCommand::with_name("client")
        .about("runs client")
        .arg(Arg::with_name("port")
            .short("p")
            .long("port")
            .value_name("CLIENT_PORT")
            .help("connect port")
            .takes_value(true)))
    .get_matches();
}
