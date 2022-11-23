# Finite State Machine

Playground for DSLRPhotoBooth project.

![Graph Visualizer](/graphviz.png)

```
digraph {
  compound=true;
  node [shape=Mrecord];
  rankdir="LR";

  ChromeWindow [label="ChromeWindow|entry / func5"];
  OffHook [label="OffHook"];
  TakePicture [label="TakePicture|entry / func2"];
  ChromeWindow -> OffHook [label="SessionEnd"];
  OffHook -> TakePicture [label="SessionStart / func1"];
  TakePicture -> TakePicture [label="CountDown"];
  TakePicture -> TakePicture [label="Printing"];
  TakePicture -> ChromeWindow [label="SessionEnd"];
  TakePicture -> TakePicture [label="SharingScreen"];
  init [label="", shape=point];
  init -> OffHook
}
```
