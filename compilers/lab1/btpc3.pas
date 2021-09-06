--- btpc64.pas	2020-02-15 14:28:10.000000000 +0300
+++ btpc64_2.pas	2021-02-13 14:37:20.350030279 +0300
@@ -795,6 +795,12 @@
   end else begin
    Error(102);
   end;
+ end else if CurrentChar='?' then begin
+   ReadChar;
+   repeat
+    ReadChar;
+   until (CurrentChar=#10) or (CurrentChar=#0);
+   GetSymbol;
  end else begin
   Error(102);
  end;
