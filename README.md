TODO:

Add a follow command. It takes a single url argument and creates a new feed follow record for the current user. It should print the name of the feed and the current user once the record is created (which the query we just made should support). You'll need a query to look up feeds by URL.

Add a GetFeedFollowsForUser query. It should return all the feed follows for a given user, and include the names of the feeds and user in the result.
Add a following command. It should print all the names of the feeds the current user is following.

Enhance the addfeed command. It should now automatically create a feed follow record for the current user when they add a feed.
