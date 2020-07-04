#!/usr/bin/perl 
use File::Basename;
use File::Copy;
use Cwd;
use utf8;

open(FH,"git diff --cached --name-status |") or die $!;

while (my $line = <FH>) {
   $line =~ /^(\w)\s+(.+)/;
   if ($1 ne "D"){
      my $f = $2;
      $f =~ tr/"//d;
      ($name,$path,$suffix) = fileparse($f,  ,qr"\..[^.]*$");
      if($suffix =~ /\.md|\.MD/) {
          my $dir = getcwd;
          print "$f\n";
          system("docker run -v $dir:/blog julianshen/ogpp '$name$suffix' > '/tmp/$name$suffix'");
          copy("/tmp/$name$suffix", "$dir/$f") or die "error moving file from /tmp/$name$suffix to $dir/$f:$!";	  
      }
   }
}
close(FH)
