# Windward

Windward is a tool to generate an etcd environment variable configuration file. It is intended to be run either during discovery or afterwards. By querying the state of the etcd cluster before attempting to join, it ensures nodes which are slow to start, or which simply fail, can still join the cluster.

Windward is also a sailing term, meaning to sail in the same direction as the wind.
