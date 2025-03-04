### 1. What is symlink and what it is doing in our context?
A symlink (symbolic link) is like a shortcut in Windows or a reference pointer in Linux. It is a special type of file that points to another file or directory. When you access the symlink, it behaves as if you're accessing the original file.

For example:

    ln -s /original/path/file.txt /shortcut/file.txt
Now, /shortcut/file.txt is just a pointer to /original/path/file.txt.

In our code, the symlink is used to bind a container's network namespace to a known location so that it can be accessed easily.

### 2. What is the purpose of symlink?
The symlink makes it easier to refer to a container’s network namespace without knowing its PID.

You can now use network tools like ip netns without manually looking up the process ID.

Example:

Before symlink:

    ip netns exec /proc/1234/ns/net ip a
After symlink:

    ip netns exec <PodName> ip a

### 3. Why is This Useful?
- Makes container networking management easier.
- Avoids the need to track changing PIDs.
- Helps in debugging and configuring container networks.

### 4. Understanding in much better way how symlink works.
A symlink (symbolic link) is like creating a shortcut to each house’s Wi-Fi router.
Instead of searching inside each house, you place an easy-to-use shortcut outside, so you can enter any house’s network quickly.

In our case:

- Every pod has its own network namespace (Wi-Fi router).
- The symlink is a shortcut to each pod’s network.
- Instead of looking deep inside the system for the pod’s network, you can quickly access it using the shortcut.

Real-World Benefit

After creating the symlink, you can easily enter a pod’s network using commands like:

    ip netns exec pod-a ip a

This lets you debug, monitor, or configure a pod’s network without entering the container itself.

#### Without Symlinks (Hard Way)
Without symlinks, you’d have to manually find the pod’s network namespace inside the system every time, which is annoying and time-consuming.

#### With Symlinks (Easy Way)
With symlinks, you just use:

    ip netns exec <pod-name> <command>
and instantly control any pod’s network like a pro! 

### 5. How do ip netns exec known the symlink path just by using pod-name?
The ip netns exec command knows the symlink path just by using the pod-name because of how Linux namespaces and the ip netns tool work together.

The ip netns command looks for /var/run/nets/<POD_NAME> symlink and gets the podname.

But, if you create a symlink in another directory such as /mysym/pod-name then, the ip netns exec command might not show required result.

For such condition, we can provide full path manually.

    ip netns exec /mysym/pod-a <command>

Or, we can use nsenter for more flexibility 

    nsenter --net=/mysym/pod-a <command>

