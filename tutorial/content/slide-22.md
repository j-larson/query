## Array comprehensions

Array comprehensions allow you to build new arrays from existing ones.

In the example on the right we build a new array containing just the
children's names for all people that have children.

<pre id="example">
	SELECT 
	    fname AS parent_name,
	    ARRAY child.fname FOR child IN tutorial.children END AS child_names
  		FROM tutorial 
		    WHERE children IS NOT NULL
</pre>
