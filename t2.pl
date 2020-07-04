#!/usr/bin/perl
use File::Basename;
use Cwd;
use utf8;
binmode(STDIN, ':encoding(utf8)');
binmode(STDOUT, ':encoding(utf8)');
binmode(STDERR, ':encoding(utf8)');

open(FH,"git diff --cached --name-status |") or die $!;
while(<FH>) {
  chomp;
  print $_;
}
