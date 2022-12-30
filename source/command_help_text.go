/*
    Meander
    A portable Fountain utility for production writing
    Copyright (C) 2022-2023 Harley Denham

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

const comm_convert = `
$1Convert Usage$0
-------------

    meander $1convert$0 input_file [output]

Convert will take a number of other screenwriting tools' 
proprietary formats and (with best practices, within reason) 
convert them to a plain-text Fountain file.

$1Supported Formats$0
-----------------

    Final Draft     $1.fdx$0
    Highland 2      $1.highland$0

Meander will check the input file extension to ensure the file 
can be handled.  This cannot be forced.`

const comm_data = `
Data generates a JSON file containing the content of and some 
additional information about a given Fountain file.

Rather than conversion to other screenplay tools, this is 
intended for use with non-screenplay software, such as 
furnishing production-tracking tools with screenplay metadata 
or dumping statistics into spreadsheets.

$1Usage$0
-----

    meander $1data$0 input_file [output] [--flags]

$1Flags$0
-----

    $1--preserve-markdown$0

    by default, the export will strip *italic*,
    +highlight+, etc. Markdown formatting characters.

    $1--revision$0

    if revision mode is enabled, a $1revised$0
    value can reflect in the content entries.

The resulting JSON blob is a dictionary containing four entries:

    + $1meta$0
    + $1title$0
    + $1characters$0
    + $1content$0

$1Meta$0
----

Meta stores information about the JSON format itself — the 
version of Meander that created it and the version of the 
format and its structures.

$1Title$0
-----

Title is a simple dictionary of the title page information.  
Keys are also sanitised: a key such as $1draft date$0 will 
become $1draft_date$0.

    "title": {
        "title": "An Movie",
        "credit": "by",
        "author": "Some Nerd",
        "draft_date": "December 2022"
    }

$1Characters$0
----------

Characters is a list of all speaking characters featured in the 
screenplay with alternate names, gender information and 
line-count based on the gender definition table.

    "characters": [
        {
            "name": "Ashby",
            "other_names": [
                "Captain Ashby",
            ],
            "gender": "male",
            "lines_spoken": 168
        },
        {
            "name": "Rosemary",
            "gender": "female",
            "lines_spoken": 220
        }
    ]

$1Content$0
-------

Content is a list of all screenplay content, tagged by type.

    "content": [
        {
            "type": "scene",
            "text": "EXT. PORCH - SUNSET",
            "scene_number": "99-A"
        },
        {
            "type": "action",
            "text": "Our heroes watch the sunset..."
        },
        {
            "type": "dialogue",
            "name": "NARRATOR",
            "dialogue": [
                "Will they find peace?",
                "(beat)",
                "Who can say..."
            ]
        },
        {
            "type": "page_break"
        }
    ]

If run with the revision mode flag, any given entry in the data 
may additionally contain a $1"revised": true$0 value.

The type table is as follows:

    scene
    dialogue
    action
    centered
    transition
    section
    synopsis
    header
    footer
    page_break
    page_number`

const comm_fountain = `
$1Fountain Cheat Sheet$0
--------------------

$1Scene Headings$0
--------------

Scene headings are lines that begin with $1INT$0 or $1EXT$0.  
All of the following words can also be used, followed directly 
by a space or a full-stop.

    INT
    EXT
    EST
    INT./EXT
    INT/EXT
    I/E
    SCENE

They are not required to be uppercase to be recognised, though 
Meander will convert them to uppercase in all default printing 
templates.

Unconventional scene headings can also be forced by leading 
them with a single full-stop.  Multiple full-stops (like an 
...ellipsis) are ignored.

    $1.$0BINOCULAR POV       


$1Scene Numbers$0
-------------

Scenes can also have manually defined scene numbers by placing 
them between two pound signs.
    
    INT. HOUSE - DAY $1#1#$0
    INT. HOUSE - DAY $1#1A#$0
    INT. HOUSE - DAY $1#1a#$0
    INT. HOUSE - DAY $1#A1#$0
    INT. HOUSE - DAY $1#I-1-A#$0
    INT. HOUSE - DAY $1#1.#$0


$1Action$0
------

An Action is what Fountain calls a simple paragraph.  It is any 
passage that isn't recognised as another element.  In Action 
elements, Fountain respects whitespace and assumes every 
carriage return is intentional.

Tabs and spaces are also retained in Action elements, allowing 
writers to indent a line.  Tabs are converted to four spaces.

All this means that writing the following (including 
indentation) would be converted exactly to the printed page:

        Scott --

        Jacob Billups
        Palace Hotel, RM 412
        1:00 pm tomorrow

You can also force an Action element can by preceding it with 
an exclamation point $1!$0


$1Character$0
---------

A Character element is any line entirely in uppercase, with one 
empty line before it and without an empty line after it.

    STEEL
    The man's a myth!

Characters can be indented for legibility in text, but most 
Fountain tools will automatically position them when printing.

Character names must begin with and include at least one 
letter, so $1A113$0 would be recognised as a valid character 
name.

You can also force a Character by preceding it with $1@$0, used 
for names that require lowercase letters or some non-Latin 
alphabets.

    $1@$0McCLANE
    Yippie ki-yay!


$1Character Extensions$0
--------------------

Character extensions are parenthetical notations that follow a 
character name on the same line.  These may be upper or 
lowercase.

    MOM (O. S.)
    Luke! Come down for supper!

    HANS (on the radio)
    What was it you said?


$1Parenthetical$0
-------------

Parentheticals follow a Character or Dialogue element, and are 
wrapped in $1(parentheses)$0.
    
    STEEL
    (starting the engine)
    So much for retirement!


$1Dialogue$0
--------

Dialogue is any text following a Character or Parenthetical 
element and ends on a linebreak.

    DAN
    Then let's retire them.
    Permanently.


$1Dual Dialogue$0
-------------

Dual, or simultaneous, dialogue is expressed by adding a caret 
^ after the second Character element.

    BRICK
    Screw retirement.

    STEEL ^
    Screw retirement.

Any number of spaces between the Character name and the caret 
are acceptable, and will be ignored. All that matters is that 
the caret is the last character on the line.


$1Lyrics$0
------

Lyrics are lines with a tilde $1~$0

    ~Willy Wonka! Willy Wonka! The amazing chocolatier!
    ~Willy Wonka! Willy Wonka! Everybody give a cheer!

Meander will style these like Dialogue (regardless of proximity 
to other dialogue elements) and italicise them.


$1Transition$0
----------

Transitions are uppercase items ending in $1TO:$0 with blank 
lines before and after:

    $1CUT TO:$0

You can also force any line to be a transition by beginning it 
with a greater-than symbol $1>$0

    $1> Burn to White$0


$1Centered Text$0
-------------

Centered text is bracketed with greater/less-than:

    $1>$0 THE END $1<$0

Spaces between the greater/less-than symbols and the target 
text are ignored.


$1Emphasis$0
--------

Fountain inherits Markdown's rules for emphasis, except that it 
reserves the use of underscores for underlining, something that 
web-focused Markdown has no need for.

    $1*$0italics$1*$0
    $1**$0bold$1**$0
    $1***$0bold italics$1***$0
    $1_$0underline$1_$0

Meander also extends the core Fountain specification.

    $1~~$0strikethrough$1~~$0
    $1+$0highlight$1+$0

You can mix and match to combine all of the above.

If you need to use any asterisks or underscores as verbatim 
text, they can be escaped with backslashes:

    $1\*$09765$1\*$0

The position of the formatting characters is also important, in 
that they should "hug" the target text and not leave "odd" 
spacing:

    This will become $1*$0italic$1*$0.
    This will not become$1*$0 italic$1*$0.

Also as with Markdown, emphasis does not extend beyond line 
breaks.  This means a single formatter does not need to be 
escaped on a given line.


$1Title Page$0
----------

The optional Title Page is always the first thing in a Fountain 
document. Information is encoded in the format $1key: value$0. 
Keys can have spaces (like $1Draft Date$0), but must end with a 
colon.

    Title:
        _**Birdman**_
        *or*
        **(The Unexpected Virtue of Ignorance)**
    Credit: by
    Author:
        Alejandro G. Iñárritu
        Nicolás Giacobone
        Alexander Dinelaris Jr.
        Armando Bó

As you can see, entering multiple lines in an entry can be done 
by indenting them beneath a key.

In Meander, the following keys are respected:

    [Center of Page]
    Title
    Credit
    Author
    Source

    [Bottom Left]
    Contact
    Copyright
    Notes

    [Bottom Right]
    Revision
    Draft Date

The Fountain specification also calls for unrecognised keys to 
be ignored, which most tools do, with the typical caveat that 
the first key must be one of the above standard ones.

Meander is no different, ignoring any trailing keys that other 
tools may implement, and indeed adds two custom ones of its own:

    Paper
    Format

These allow the user to simply specify their choice of template 
and paper size without needing to always enter it in the 
command parameters.


$1Page Breaks$0
-----------

Page Breaks are indicated by a line containing three or more 
consecutive equals signs, and nothing more.

    $1===$0


$1Punctuation$0
-----------

The core Fountain spec doesn't do any of that, because the 
screenplay typographical convention is to emulate a typewriter. 
However you type your apostrophes, quotes, dashes, and dots, 
that's how they'll wind up in the screenplay.


$1Indentation$0
-----------

Tabbed or spaced indentation is ignored by Fountain and treated 
as if they were not there.

This means you can manually indent your character dialogue, 
should you choose to do so, and Meander will have no trouble 
reading it.

                SOME DUDE
            (angry)
        Hey!  That's my waffle!

The exception to this is in Action elements, as aforementioned, 
where indentation is respected.


$1Meander's Syntax Extensions$0
---------------------------

Meander extends the core Fountain syntax with some of its own 
features and some borrowed from other screenwriting tools.


$1Directives$0
----------

Meander recognises Highland 2-style directives.

$1Includes$0

    {{include: scenes/1_02.fountain}}

Includes let you...

$1Counters$0

    {{series}}
    {{figure}}
    {{chapter}}
    {{panel}}

Counters are substitute values that increment automatically 
every time they are encountered.  You can set (or restart) the 
starting value to any number like so:

    {{series: 1}}

$1Timestamps$0

    {{timestamp: dd MM yyyy}}

Timestamps will substitute the current date / time, helpful for 
those of us who forget to change the dates on their drafts.  
Timestamp uses a sensible default, but you can use a custom 
pattern using any of the following characters:

    $1Days$0
    d       2
    dd      02
    E       Mon
    EEEE    Monday

    $1Months$0
    M       1
    MM      01
    MMM     Jan
    MMMM    January

    $1Years$0
    y       2022
    yy      22

    $1Hours$0
    h       3
    hh      03
    HH      15
    a       PM

    $1Minutes$0
    m       4
    mm      04

    $1Seconds$0
    s       5
    ss      05

    $1Milliseconds$0
    SSS     .000`

const comm_gender = `
$1Gender Usage$0
------------

    meander $1gender$0 [input.fountain]

Gender performs simple analysis of your characters' gender 
identities, providing a detailed print-out of how they break 
down across a script.

The Gender command needs to be given data to ensure it can 
provide accurate statistics.  It expects this in the form of a 
boneyard comment:

/*
    [gender.male]
    Ashby
    Jenks

    [gender.female]
    Rosemary
    Sissix
    Kizzy

    [gender.<custom>]
    Dr. Chef

    [gender.ignore]
    Market Stall Owner
*/

This special boneyard must be found somewhere in the input 
file.  If there is more than one, the first one is used.  The 
first non-whitespace text inside it must be a [gender] heading.

Indentation and casing of headings and characters are ignored, 
but the table must be arranged line-by-line.  Blank lines are 
skipped.

There are no reserved terms for gender identities.  All terms 
used are technically "custom", so any number of different 
identities can be specified, useful for non-binary and queer 
characters, as well as for instances like non-humans in science 
fiction.

The only exception to this is "ignore".  Characters assigned to 
"ignore" will not be counted in the statistics, useful for 
single-appearance characters or special cases.

Any characters found in the screenplay but *not* in the gender 
table will be reported as "unknown" and classified in the 
statistics under that additional group.

$1Names$0
-----

Characters can also have multiple names throughout the script, 
for a variety of reasons.  Meander can handle this if provided 
with all the variants, separated by pipes:

/*
    [gender.male]
    Ashby | Captain Ashby
*/

The first name in the entry is used as their canonical name for 
all subsequent output.`

const comm_help = `
$1Usage$0
-----

    meander [command] [input.fountain] [output.pdf] [--flags]

Meander is designed to run as sensibly as possible with as few 
options as possible:

    + If no command is supplied, Meander will
      render a PDF
    + If no input is supplied, Meander will look
      for a $1.fountain$0 (or $1.ftn$0) file, giving
      priority to ones named $1root$0, $1main$0 or $1manifest$0.
    + If no output is supplied, Meander will swap
      the extension to an appropriate one for the
      chosen command.

All this means you can simply call Meander (with no arguments) 
in the same location as a single-file screenplay, making 
generic build systems for text editors extremely easy to 
configure.

$1Commands$0
--------

    $1render$0    render input file to PDF (default)
    $1gender$0    display gender analysis statistics
    $1merge$0     merge a multi-file document
    $1convert$0   (experimental) convert from other software
    $1help$0      print this message and others

$1Help$0
----

Use $1meander help [command]$0 for more information on the 
above commands, but also see the additional help topics 
available below:

    $1fountain$0  fountain cheat sheet`

const comm_merge = `
$1Merge Usage$0
-----------

    meander $1merge$0 input_file [output]

Merge collapses a multi-file screenplay with include directives 
into a single text file.

If the output is unspecified, Meander will save the result with 
a suffix to prevent clobbering the original file.  You can 
forcibly overwrite the original by manually specifing it as 
both arguments.

$1Include Directive$0
-----------------

In Highland-flavoured Fountain, an include directive allows 
linking to the contents of another file.  This allows writing a 
film with each scene in a separate file, or a manuscript by 
chapters.

Includes take the form:

    $1{{include: some/file.fountain}}$0

The path to the included file should be relative to the file in 
which the directive is written.

They can be specified anywhere in text and can be nested in 
multiple layers of children.  Critically, because of this 
nested infinity, there is no cycle-safety.  Meander will get 
stuck if a loop between included files is created.`

const comm_render = `
Render will take a Fountain file and convert it to a fully 
formatted PDF document.

$1Render Usage$0
-----

    meander $1render$0 [input.fountain] [output.pdf] [--flags]

The first option of any flag listed below is also the default, 
which means it does not need to be specified unless performing 
an override.

$1Scene Numbers$0
-------------

    $1--scene -s$0

    input           uses numbers from input Fountain text
    remove          removes all scene numbers from output
    generate        creates new, sequential scene numbers

$1Format / Layout$0
---------------

    $1--format -f$0     (or title page) $1format: manuscript$0

    screenplay      standard film script or screenplay
    stageplay       stage directions and right-aligned dialogue
    manuscript      standard wide-spaced novel manuscript
    graphicnovel    sections added for panel directions

$1Paper Size$0
----------

    $1--paper -p$0      (or title page) $1paper: A4$0

    A4
    USLetter

Note that for maximum compatibility, "paper" and "format" 
should *not* be the first entries in the title page.  Most 
parsers will reject the entire title page if the first entry is 
unknown to them, but will quietly skip later ones.  This is 
never a guarantee, but it may be useful.

$1Force Hidden Syntaxes$0
---------------------

    $1--notes$0
    $1--synopses$0
    $1--sections$0`