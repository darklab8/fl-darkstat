###############################################################################
#FLBini.pm
#Last Change: 2005-04-14
#Copyright (c) 2004 Marc-Seabstian "Maluku" Lucksch
#Version 0.21
####################
#This file is part of the Plasma project, a parser library for an all-purpose
#ASCII file format. More information can be found on the project web site
#at http://aurora.dbb.at/ .
#
#FLBini.pm is published under the terms of the MIT license, which basically 
#means "Do with it whatever you want". For more information, see the license.txt
#file that should be enclosed with plasma distributions. A copy of the license
#is (at the time of this writing) also available at
#http://www.opensource.org/licenses/mit-license.php .
###############################################################################

package Plasma::FLUtf;
use strict;
use warnings;
use Data::Dumper;
use Tie::InsertOrderHash;
use Carp;
no warnings 'recursion';

=head1 NAME

FLUtf - Perl extension for working with UTF Files

=head1 Synopsis

	use FLUtf;
	open FILE,"crypt.cmp"; #or .utf, .3db, .txm, .mat, .ale, .vms, .dfm or maybe some more
	binmode FILE;
	my $crypted = do {local $/; <FILE>};
	close FILE;
	my $tree=UTFread($crypted);
	$code = UTFwriteUTF($tree);
	open FILE, ">out.cmp"
	binmode FILE;
	print FILE, $code;
	close FILE;


=head1 DESCRIPTION

This Module provides the ability to decode utf files for the Mircrosoft game "Freelancer"
Those are named UTF files because of their header

=head1 WARNING

The read routines return a tied InsertOrderHash (get it from CPAN its pure perl), so just keep the reference you got and work with it instead of using

	my %mysuperhash = %{UTFread($crypted)}
	%mysuperhash{NewEntry} = "Data"
	# Now the order is destroyed
	# This is better:
	my $tree=UTFread($crypted);
	$tree->{NewEntry} = "Data"
	

I have no idea how important the order of the elements is for Freelancer, but better keep it this way.
Of course all subhashes are also tied
	#Bad code example
	$tree->{copyme}={%{$tree}};
	#Good code:
	tie my %newhash,'Tie::InsertOrderHash';
	%newhash=(%{$tree}) #Copy tree, preserve order #At least last time I tested
	$tree->{copyme}=\%newhash;
	#not this:
	$tree->{copyme}={%newhash}; #Looses tiedness too.



=head1 FUNCTIONS

=cut

# Perl-Port by Maluku (fl@maluku.de)
our @ISA=qw/Exporter/;
our @EXPORT=qw/UTFread UTFwriteUTF/;
our @EXPORT_OK=qw/UTFread UTFwriteUTF/;

my $datas;
my $strings;
my $pointer;
my %strings;
my %datas;
my %offsets;
=head2 $chars = get ($string,$offset,$amount);

returns $amount (or less) chars of $string starting from $offset.
Also increases $offset.

=cut


sub get {
	my $off=$_[1];
	$_[1]+=$_[2];
	return substr($_[0],$off,$_[2]);
}

=head2 $string = string ($offset);

Returns a string of $stringlib, which is a compilation of \0 strings starting a $offset.
Returns string without the trailing \0.

=cut

sub string {
	return "" if $_[0] >= length $strings; #THIS IS AN ERROR IN THE UTF FILE
	my $shift=index($strings,"\0",$_[0]) - $_[0];
	return substr($strings,$_[0],$shift);
}
sub data {
	return "" if $_[0] >= length $datas; #THIS IS AN ERROR IN THE UTF FILE
	return substr($datas,$_[0],$_[1]);
}


=head2 $tree = UTFread ($crypted);

Reads a tree (Hash of Hashes of H...)


=cut

sub UTFread {
	my $d = UTFreadUTF(@_);
	#print Data::Dumper->Dump([$d]);
	return $d;
}


#node
#{
#	dword sibling_offset
#	dword string_offset
#	dword flags
#	dword zero (seems to be always zero, meaning unknown)
#	dword child_offset
#	dword allocated_size
#	dword size1
#	dword size2
#	dword time1
#	dword time2
#	dword time3
#} = 44 bytes

sub UTFreadUTFrek {
	local $_;
	my $tree=shift;
	my $i=shift;
	return {} if $offsets{$i}++;
	return {} if $i > length($tree)-44;
	tie my %data => 'Tie::InsertOrderHash';
	#print $i;
	
	my ($silb, $name, $flags, $z, $childoffset, $alloc, $size, $size2, $time1, $time2, $time3) = unpack("VVVVVVVVVVV", substr($tree,$i));
	#print ' $silb, $name, $flags, $z, $childoffset, $alloc, $size, $size2, $time1, $time2, $time3 = ',"($silb, $name, $flags, $z, $childoffset, $alloc, $size, $size2, $time1, $time2, $time3)\n";
	# Now we must use all the stuff here or warning will annoy us:
	($time1, $time2, $time3) = ($time1, $time2, $time3);
	$size=$size2 if ($size2 < $size);
#	die "Allocation error in $i" if ($alloc < $size);
	if ($flags & 0x10 and not $flags & 0x80) {
		$data{string($name)}=UTFreadUTFrek($tree,$childoffset);
	}
	else {
		$data{string($name)}=data($childoffset,$size);
	}
	if ($silb) {
		my $data=UTFreadUTFrek($tree,$silb);
		#print $data,"->$i-->$silb\n";
		%data=(%data,%{$data});
	}
	return \%data;

}


sub packString {
	my $string = shift;
	if (exists $strings{$string}) {
		return $strings{$string};
	}
	else {
		$strings{$string} = length($strings);
		$strings.=$string."\0";
		return $strings{$string};
	}

}

sub packData {
	my $data = shift;
	if (exists $datas{$data}) {
		return $datas{$data};
	}
	else {
		$datas{$data} = length($datas);
		$datas.=$data."\0";
		return $datas{$data};
	}

}

sub UTFwriteUTFrek {
	local $_;
	my $tree=shift;
	my $name=shift;
	my $silb=shift;
	$pointer+=44;
	my $start=$pointer;
	die "Can't pack other then scalar or Hashrefs" if ref $tree and ref $tree ne "HASH";
	return pack ("VVVVVVVVVVV",$silb?$pointer:0, packString($name), 0x80, 0, packData($tree), length($tree), length($tree), length($tree), 0, 0, 0) unless ref $tree;
	my $code="";
	my @list = keys(%$tree);
	#print "name = $name, pointer = $pointer, data = ".(ref ($tree) || "scalar")."\n";
	foreach (0 .. $#list) {
		$code.=UTFwriteUTFrek($tree->{$list[$_]},$list[$_],($_!=$#list));
	}
	#print "name = $name, pointer = $pointer, data = ".(ref ($tree) || "scalar")."\n";
	return pack ("VVVVVVVVVVV",$silb?$pointer:0, packString($name), 0x10, 0, $start, 0,0,0,0,0,0).$code;
}


#header
#{
#	dword "UTF "
#	dword 0x101
#	dword tree_segment_offset
#	dword size_of_tree_segment
#	dword header_offset? (0) ##I think its here also a First element of the treeoffset
#	dword size_of_header (44) ##I think it is more a size of entry
#	dword string_segment_offset
#	dword space_allocated_for_string_segment
#	dword size_of_string_segment_actually_used
#	dword data_segment_offset
#	dword unknown (seems to be zero most of the times) #Possible first entry of data segment (after deletion of an entry)
#} = 44 bytes

sub UTFreadUTF{
	my $code=shift;
	my $i=0;
	%offsets = ();
	if (substr($code,$i,4) eq "UTF ") {
		$i+=4;
		my ($ver,$treeoffset,$treesize,$treefirst,$treeelemsize,$stringoffset,$stringspace,$stringsize,$dataoffset,$datafirst)=unpack("VVVVVVVVVV",get($code,$i,40));
		# We don't use those now, not sure what they are for anyway
		$datafirst = 0 unless $datafirst;
		$treefirst = 0 unless $treefirst;
		$treeelemsize = 0 unless $treeelemsize;
		# We don't need this one either, do we?
		$stringsize = 0 unless $stringsize;
		# Splitting the parts
		$strings=substr($code,$stringoffset,$stringspace);
		$datas=substr($code,$dataoffset);
		my $tree=substr($code,$treeoffset,$treesize);
		return UTFreadUTFrek($tree,0);


	}
	else {
		die "NOT a UTF File";
	}
	
}

sub UTFwriteUTF{
	my $tree=shift;
	my $i=0;
	$strings="";
	$datas="";
	my $code = "";
	%strings = ();
	%datas = ();
	$pointer=0;
	my @list = keys(%$tree);
	foreach (0 .. $#list) {
		$code.=UTFwriteUTFrek($tree->{$list[$_]},$list[$_],($_!=$#list));
	}
	my $string=$strings;
	$string.="\0" for(length($strings) .. (int(length($strings)/32)+1)*32); #Just some fun stuff.
	return "UTF ".pack("VVVVVVVVVV",0x101,44+12,length($code),0,44,44+12+length($code),length($string),length($strings),44+12+length($code)+length($string),0)."000000000000".$code.$string.$datas;
	
}


=cut

=head1 AUTHOR

Marc "Maluku" Sebastian Lucksch
plasma@marc-s.de

=cut
1;
