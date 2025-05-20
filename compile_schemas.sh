#!/bin/bash
flatc --go -o client/flatbuffers schemas/*.fbs
flatc --go -o server/flatbuffers schemas/*.fbs
# flatc --go -o schemas/gen schemas/*.fbs
