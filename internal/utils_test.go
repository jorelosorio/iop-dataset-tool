package internal

import (
	"math"
	"reflect"
	"testing"
)

const text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas auctor lacus eget maximus aliquam. Phasellus vitae enim urna. Etiam ultrices mauris non leo dignissim interdum. Praesent suscipit ultricies imperdiet. Donec mattis egestas bibendum. Aliquam justo dui, fermentum eget consequat a, faucibus non lacus. Quisque ac tortor scelerisque, tempor augue at, aliquam lorem. Vestibulum non placerat tellus. Sed viverra cursus erat sed cursus. In efficitur rutrum felis sit amet rutrum. Praesent pellentesque purus sit amet ligula pellentesque consequat. Suspendisse leo sem, finibus scelerisque auctor interdum, dapibus ac sapien. Etiam massa metus, gravida et nisi vitae, volutpat cursus sapien. Aenean pharetra vulputate nisl. Nullam placerat pretium risus, vel gravida mi interdum ut. Aliquam erat volutpat.
Cras vel elit posuere tellus egestas gravida. Nulla suscipit orci ac egestas laoreet. Ut diam magna, pulvinar eget pharetra vel, scelerisque vitae augue. Duis rutrum ut nunc ut venenatis. Pellentesque nibh turpis, ultricies a arcu ac, congue tempor est. Praesent venenatis lectus eget mi mattis venenatis. Duis a metus fermentum, semper nulla eget, mollis sapien. Suspendisse fringilla tempor nulla, ornare imperdiet augue consectetur sed. Praesent libero mi, aliquet sit amet maximus vitae, placerat at arcu. Maecenas suscipit, ex ut tempor dignissim, est magna fringilla ex, non rutrum magna erat egestas orci. Duis tempor varius ligula sed fermentum. In feugiat, mauris eget fringilla gravida, lacus nisl convallis velit, vitae imperdiet nisi tellus eget diam. Nulla facilisi.
Nullam et dignissim leo. Duis mattis odio non lorem mollis, id venenatis quam efficitur. Donec in varius mauris. Integer ut tortor auctor risus vehicula iaculis at non erat. Aliquam placerat mattis ullamcorper. Interdum et malesuada fames ac ante ipsum primis in faucibus. Quisque tincidunt, leo ac interdum rhoncus, enim ipsum auctor erat, sagittis feugiat massa ligula porta nunc. Phasellus dolor metus, mollis eget arcu commodo, elementum scelerisque eros.
Sed porta mauris vitae molestie accumsan. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Morbi ultricies vitae erat sed vehicula. Vivamus euismod justo vitae enim fringilla, luctus rutrum erat aliquam. Aenean venenatis, magna vitae facilisis vehicula, turpis odio commodo mauris, vitae varius libero magna non mauris. Cras lectus turpis, pellentesque quis justo ac, gravida blandit dolor. Donec vitae fringilla ligula.
Praesent at augue vel lorem porttitor blandit. Vestibulum eget accumsan felis. In in mauris accumsan, volutpat enim vel, semper metus. Suspendisse sit amet mi sed orci rhoncus sagittis. Sed blandit diam nunc, nec tincidunt risus tempus quis. Donec nec felis.`

func TestSplitTextIntoChunks(t *testing.T) {
	chunkSize := int(math.Round(float64(len(text) / 3)))
	chunks := SplitTextIntoChunks(text, chunkSize)
	if !reflect.DeepEqual(3, len(chunks)) {
		t.Errorf("Expected %v, got %v", chunkSize, len(chunks))
	}
}

func TestSplitTextWithSmallChunkSize(t *testing.T) {
	chunks := SplitTextIntoChunks(text, 56)
	smallText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	if !reflect.DeepEqual(smallText, chunks[0]) {
		t.Errorf("Expected %v, got %v", smallText, chunks[0])
	}
}
