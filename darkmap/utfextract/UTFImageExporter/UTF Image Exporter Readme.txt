
				UTF Image Exporter v0.0.2
				
	Copyright 2010, Fenris_Wolf, http://yspstudios.linuxniche.net

-------------------------------------------------------------------------
	
	Mass exporting (and importing) of .tga and .dds files from Freelancer UTF 
	files (.cmp, .3db, .mat and .txm).
	
	Exporting takes all UTF files in a folder,  and exports the images into a 
	seperate output folder as: utf_filename.cmp\image.tga
	
	Importing takes the same directory structure (utf_filename.cmp\image.tga)
	and sticks them back in the files in the UTF folder.

	
	Limits:
	only imports images if they already exist in the UTF. Doesnt create new 
	branches on the UTF tree.
	
	when importing if the image in the UTF is a .dds (its name in the UTF tree
	would still have the .tga extension), the file getting imported must also be 
	a .dds (but with a .dds extension, NOT .tga)

	Some files have odd names in UTF files: somefile.TGA1231243 for these
	exports and imports use the names in the UTF - no proper .tga or .dds 
	file extension (you will have to manually do these)
	
	The results.log file is wiped every time the program starts. If you which to
	archive this after running, you must do it manually.
	
	If theres over 250 UTF files to be scanned, Verbose mode is disabled.
	This only effects the output window, everything is stilled logged to 
	the results.log



	Credits to maluku for the utf perl module
	
	changelog:
	
	v0.0.2
	.dfm file support
	added multiple MIP#  per filename.tga (MIP0 MIP1, etc)
	added recursive UTF file paths ^^
	added log file
	
	
