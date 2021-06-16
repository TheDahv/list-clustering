# List Clustering

This project explores the task of identifying groups of related items in a list.
Typical applications might include determining "similarity" with respect to some
metadata about an item in a list.

To put it more concretely, it could also allow you to surface clusters of
related keywords by looking at the search results they trigger in a given search
engine.

It leverages 2 techniques to make this work:

- Rank biased overlap: an algorithm that determines the similarity of two ranked
  lists
- Markov clustering: an algorithm to determine "neighborhoods" of related nodes
  within a graph

## TODO

- [x] Implement RBO for an interface that a given struct can implement to
  surface the list data examined to determine similarity
- [ ] Implement Markov clustering
  - [jamesneve/go-markov-cluster](https://github.com/jamesneve/go-markov-cluster)
    is the closest candidate, but uses old versions of gonum and isn't as
    flexible as I'd like
- [ ] Look at ways to split the RBO computation; compute RBO for a subset and
  cluster, then compare the remainder only to cluster representatives. Matches
  above a treshold join the cluster, and others begin new clusters
