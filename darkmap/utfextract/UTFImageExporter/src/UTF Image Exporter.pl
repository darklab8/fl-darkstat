use lib '../lib';
use Plasma::FLUtf;
use Tie::InsertOrderHash;
use Win32::GUI;
use strict;
use warnings;
use Cwd;

use File::Find;
use File::Path;

########################
# GUI
my $MainWindow;
my $ModelPathTextfield;
my $ImagePathTextfield;
my $ExportRadio;
my $ImportRadio;
my $ResultsTextfield;
my $RecursiveCheck;
my $VerboseCheck;

########################
my $ModelPath;
my $ImagePath;
my $RunMode;
my $RecursiveMode;
my $Verbose = 1;

my @Status;

open LOGFILE, '>', 'results.log';
close LOGFILE;
########################
sub logLine {
	my $line = shift;
	print LOGFILE "$line\n";
	if ($Verbose || $_[0]) {
		$ResultsTextfield->Append( $line."\r\n");
	}
}

sub readModel { # intial UTF read
	open FILE, $_[0];
	binmode FILE;
	my $crypted = do {local $/; <FILE>};
	close FILE;
	my $tree = UTFread $crypted;
	return $tree;
}

sub getLibrary { # reference to our texture library
	my $file = shift;
	my $tree = readModel "$ModelPath/$file";
	my $root = $tree->{'\\'};
	my $obj;
	while (my ($k, $v) = each %{$root} ) {
		$obj = $v and last if $k =~ /^texture library$/i;
	}
	return (undef, $tree) unless $obj;
	# yes, i'm cleaning a reference
	$obj =~ s/^\s+//g; 
	$obj =~ s/\s+$//g;
	$obj =~ s/\000//g; 
	$Status[1]++; # inc one text lib
	return ($obj, $tree);
}

sub getUTFFiles { # 

	my @return;
	unless ($RecursiveMode) { 
		opendir DIR, $ModelPath;
		my @files = readdir DIR;
		closedir DIR;
		foreach my $file (@files) {
			next if -d $file;
			push @return, $file if $file =~ /\.(3db|cmp|mat|txm|dfm)$/i;
		}
		return @return;
	}

	find sub{
		my $file = $File::Find::name;
		return undef unless $file =~ /\.(3db|cmp|mat|txm|dfm)$/i;
		$file =~ s/^\Q$ModelPath\E\///i;
		push @return, $file;
	}, $ModelPath;
	
	return @return;
}

sub saveUTF {
	my ($file, $tree) = (shift, shift);
	logLine "Updating $file";

	my $code = UTFwriteUTF $tree;
	open FILE, '>', "$ModelPath/$file";
	binmode FILE;
	print FILE $code;
	close FILE;
	$Status[2]++;
}

sub saveImage { # saves a exported image
	my ($file, $data) = (shift, shift);
	logLine "Exporting $file";

	open FILE, '>', "$ImagePath/$file" or return undef;
	binmode FILE;
	print FILE $$data;
	close FILE;
	$Status[3]++;
}

sub readImage {
		my $file = shift;
		open FILE, "$ImagePath/$file" or return undef;
		binmode FILE;
		my @d = <FILE>;
		close FILE;
		
		logLine "Importing $file";
		return join( '', @d);
}

sub exportImages { # all images in a UTF
	my $filename = shift;
	my ($obj) = getLibrary $filename;
	return undef unless $obj;

	mkpath "$ImagePath/$filename";

	while (my ($kfile, $v) = each %$obj) {
		my @vkeys = keys %$v;
		next unless @vkeys;

		if (@vkeys == 1) { # normal use only 1 MIP0 or MIPS
			my $data = \$v->{  $vkeys[0] }; 
			next unless $$data;
			if (substr( $$data,0, 3) eq 'DDS') {
				$kfile =~ s/tga$/dds/i;
			}
			 saveImage "$filename/$kfile", $data;
			 next;
		}
		
		# here we have more then 1 key, multiple images, .dfm file?
		mkdir "$ImagePath/$filename/$kfile";
		foreach my $key (@vkeys) {
			my $data = \$v->{  $key }; 
			next unless $$data;
			
			if (substr( $$data, 0, 3) eq 'DDS') {
				saveImage "$filename/$kfile/$key.dds", $data;
			}
			else {
				saveImage "$filename/$kfile/$key.tga", $data;
			}
		}
	}
}

# need 2 import modes..a 'safe' mode (images must already exist), and a add mode

sub importImages {
	my $filename = shift;
	return undef unless -d "$ImagePath/$filename"; # no import folder
	my ($obj, $tree) = getLibrary $filename or return undef;
	
	my $changed;
	while (my ($kfile, $v) = each %$obj) {
		my @vkeys = keys %$v;
		next unless @vkeys;
		if (@vkeys == 1) { # normal use only 1 MIP0 or MIPS
			my $data = \$v->{  $vkeys[0] }; 
			next unless $$data;
			if (substr( $$data,0, 3) eq 'DDS') {
				$kfile =~ s/tga$/dds/i;
			}

			$data = readImage "$filename/$kfile";
			next unless $data;
			$changed = 1;
			$Status[3]++;
			$v->{ $vkeys[0] } = $data;
			next;
		}
		
		# here we have more then 1 key, multiple images, .dfm file?
		foreach my $key (@vkeys) {
			my $data = \$v->{  $key }; 
			next unless $$data;
			
			if (substr( $$data, 0, 3) eq 'DDS') {
				$data = readImage "$filename/$kfile/$key.dds";
			}
			else {
				$data = readImage "$filename/$kfile/$key.tga";
			}
			next unless $data;
			$changed = 1;
			$Status[3]++;
			$v->{ $key } = $data;
		}		
	}
	saveUTF $filename, $tree if $changed;
	
}

sub checkPaths {
	unless ($ModelPath) {
		$ResultsTextfield->Text( "ERROR! UTF directory not specified" );
		return undef;
	}
	unless ($ImagePath) {
		$ResultsTextfield->Text( "ERROR! Image directory not specified" );
		return undef;
	}
	unless (-d $ModelPath) {
		$ResultsTextfield->Text( "ERROR! Cant Find UTF directory at: $ModelPath" );
		return undef;
	}
	unless (-d $ImagePath) {
		$ResultsTextfield->Text( "ERROR! Cant Find Image directory at: $ImagePath" );
		return undef;
	}

	if (lc $ModelPath eq lc $ImagePath) {
		$ResultsTextfield->Text( "ERROR! UTF and Image directories cannot be the same" );
		return undef;
	}
	return 1;
}

sub resetStatus {

	$ResultsTextfield->Text( "" );
	$ModelPath = $ModelPathTextfield->Text();
	$ImagePath = $ImagePathTextfield->Text();
	
	$RunMode = $ExportRadio->Checked();
	$RecursiveMode = $RecursiveCheck->Checked();
	$Verbose = $VerboseCheck->Checked();
	
	$ModelPath =  cwd().'/'. $ModelPath unless $ModelPath =~ /^\S\:/;
	$ModelPath =~ s/\\/\//g; # replace any \ with /

	$ImagePath =  cwd().'/'. $ImagePath unless $ImagePath =~ /^\S\:/;
	$ImagePath =~ s/\\/\//g; # replace any \ with /
	
	@Status = ( 
		0, # 0 - total UTF files in UTF dir
		0, # 1 - total UTF files with texture libraries
		0, # 2 - total UTF files modified
		0, # 3 - total Images exported / imported
	);
}

sub runIt {

	resetStatus;
	
	return undef unless checkPaths;
	my @files = getUTFFiles;

	$Status[0] = scalar @files;
	if ($Status[0] > 250) {
		return undef if Win32::GUI::MessageBox( 0, $Status[0]." UTF Files found, are you sure? This may take a while, verbose will be disabled", "Warning:", 0x0001 ) == 2;
		$Verbose = 0;
	}

	my $rmode = 'Import' ;
	$rmode = 'Export' if $RunMode;

	open LOGFILE, '>>', 'results.log';

	######################################
	# Print some header info
	my $time = localtime();
	logLine "$rmode Mode - $time", 1;
	
	logLine "UTF Folder: $ModelPath", 1;
	logLine "Image Folder: $ImagePath", 1;
	
	######################################
	# Start the check
	my $start_time = time();
	foreach my $fi (@files) {
		next unless $fi;
		logLine "Checking $fi";
		if ($RunMode) {
			exportImages $fi;
		}
		else {
			importImages $fi;
		}
	}
	
	######################################
	# done
	my $end_time = time();
	logLine "--------------------------------------------------------------------------------", 1;
	logLine "Total UTF Files Read: $Status[0]", 1;
	logLine "Total UTF Texture Libraries: $Status[1]", 1;
	logLine "Total UTF Files Modified: $Status[2]", 1;
	logLine "Total Images ".$rmode. "ed: $Status[3]", 1;
	logLine "Total Time: ".($end_time - $start_time)." seconds", 1;

	logLine '';
	close LOGFILE;

}

sub buildMainWindow {
    $MainWindow = Win32::GUI::Window->new(
        -name => 'MainWindow',
        -text => 'UTF Image Exporter',
        -width => 500,
        -height => 310
	);
	$ModelPathTextfield = Win32::GUI::Textfield->new(
        $MainWindow,
		-name		=> 'ModelPathTextfield',
        -pos        => [10, 10],
		-size		=> [ 300, 20 ]
	);
    Win32::GUI::Button->new(
        $MainWindow,
        -name       => "ModelPathButton",
        -pos        => [320, 10],
        -size       => [ 140, 20 ],
        -text       => 'UTF Directory Path'
    );
	$ImagePathTextfield = Win32::GUI::Textfield->new(
        $MainWindow,
		-name		=> 'ImagePathTextfield',
        -pos        => [10, 40],
		-size		=> [ 300, 20 ]
	);
    Win32::GUI::Button->new(
        $MainWindow,
        -name       => "ImagePathButton",
        -pos        => [320, 40],
        -size       => [ 140, 20 ],
        -text       => 'Image Directory Path'
    );

	##########################################
    $ExportRadio = Win32::GUI::RadioButton->new(
        $MainWindow,
        -pos        => [320, 70],
        -text       => 'Export Images from UTF',
		-checked	=> 1
    );
	
    $ImportRadio = Win32::GUI::RadioButton->new(
        $MainWindow,
        -pos        => [320, 90],
        -text       => 'Import Images into UTF',
    );

	$RecursiveCheck = Win32::GUI::Checkbox->new(
        $MainWindow,
        -pos        => [320, 110],
        -text       => 'Recursive Mode',
	);
	
	$VerboseCheck = Win32::GUI::Checkbox->new(
        $MainWindow,
        -pos        => [320, 130],
        -text       => 'Verbose',
		-checked => 1
	);


    Win32::GUI::Button->new(
        $MainWindow,
        -name       => "StartButton",
        -pos        => [320, 230],
        -size       => [ 140, 40 ],
        -text       => 'Start'
    );
	
	##########################################
	$ResultsTextfield = Win32::GUI::Textfield->new(
        $MainWindow,
		-name		=> 'ResultsTextfield',
        -pos        => [10, 70],
		-multiline => 1,
		-readonly => 1,
		-autohscroll	=> 1,
		-size		=> [ 300, 200 ]
	);

}
############################################################################################

buildMainWindow;
$MainWindow->Show();


#$ModelPathTextfield->Text('C:\\Games\\Freelancer\\Data');
#$ModelPathTextfield->Text('utf_temp');
#$ImagePathTextfield->Text('img_temp2');

Win32::GUI::Dialog();

############################################################################################
#       GUI EVENT HANDLERS


sub MainWindow_Terminate {
	close LOGFILE;
    -1;
}
 sub ModelPathButton_Click {
	my $dir = Win32::GUI::BrowseForFolder(
		-title => 'Select folder containing UTF files',
		-folderonly => 1);
	$ModelPathTextfield->Text($dir) if $dir;
}

sub ImagePathButton_Click {
	my $dir = Win32::GUI::BrowseForFolder(
		-title => 'Select root folder containing image files',
		-folderonly => 1);
	$ImagePathTextfield->Text($dir) if $dir;
}
sub StartButton_Click {
	runIt;
	1;
}

1;
__END__
