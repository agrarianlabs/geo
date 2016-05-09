package geo

import "testing"

func TestPolygonMarshal(t *testing.T) {
	p := &Polygon{
		{1.2, 3.4},
		{5.6, 7.8},
	}
	expected := `{"type":"Polygon","coordinates":[[1.2,3.4],[5.6,7.8]]}`
	got, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, string(got))
	}
}

func TestPolygonScan(t *testing.T) {
	// Good
	for _, testcase := range []struct {
		WKT      string
		Expected Polygon
	}{
		{
			WKT: "POLYGON((1.2 3.4, 5.6 7.8, 6.2 1.5, 1.2 3.4)",
			Expected: Polygon{
				{1.2, 3.4},
				{5.6, 7.8},
				{6.2, 1.5},
				{1.2, 3.4},
			},
		},
	} {
		p := &Polygon{}
		if err := p.Scan(testcase.WKT); err != nil {
			t.Fatal(err)
		}
		for i, coord := range testcase.Expected {
			if expected, got := coord[0], (*p)[i][0]; expected != got {
				t.Fatalf("expected %f, got %f", expected, got)
			}
			if expected, got := coord[1], (*p)[i][1]; expected != got {
				t.Fatalf("expected %f, got %f", expected, got)
			}
		}
	}
	// Bad
	for _, testcase := range []interface{}{
		"POLYGON((1.2, 3.4, 5.6, 7.8))",
		[]byte("POLYGON((1.2, 3.4, 5.6, 7.8))"),
		7,
		"POLYGON(1.2 3.4 5.6 7.8)",
		"POLYGON((1.2 3.4 5.6 7.8)}",
		"PIKACHU",
	} {
		p := &Polygon{}
		if err := p.Scan(testcase); err == nil {
			t.Fatalf("expected err, got nil")
		}
	}
}

func TestPolygonValue(t *testing.T) {
	var (
		p = Polygon{
			{1.2, 3.4},
			{5.6, 7.8},
			{8.7, 6.5},
			{4.3, 2.1},
		}
		expected = `POLYGON((1.2 3.4, 5.6 7.8, 8.7 6.5, 4.3 2.1))`
	)
	value, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("expected string, got %T", value)
	}
	if expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestPolygonEmpty(t *testing.T) {
	var (
		p        = Polygon{}
		expected = "POLYGON EMPTY"
	)
	value, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("expected string, got %T", value)
	}
	if expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestPolygonContains(t *testing.T) {
	for _, testcase := range []struct {
		Poly    Polygon
		Inside  []Point
		Outside []Point
	}{
		// Square
		{
			Poly: Polygon{
				{0, 0},
				{2, 0},
				{2, 2},
				{0, 2},
				{0, 0},
			},
			Inside: []Point{
				{1, 1},
			},
			Outside: []Point{
				{4, 1},
			},
		},
		// Hexagon
		{
			Poly: Polygon{
				{0, 1},
				{1, 2},
				{2, 1},
				{2, 0},
				{1, -1},
				{0, 0},
				{0, 1},
			},
			Inside: []Point{
				{1, 0},
			},
		},
		// A tilted quadrilateral
		{
			Poly: Polygon{
				{-1, 10},
				{10, 1},
				{1, -10},
				{-10, -1},
				{-1, 10},
			},
			Inside: []Point{
				{2, 2},
				{2, -2},
			},
		},
	} {
		if testcase.Inside != nil {
			for _, point := range testcase.Inside {
				if !testcase.Poly.Contains(point) {
					t.Fatalf("Expected polygon %v to contain point %v", testcase.Poly, point)
				}
			}
		}
		if testcase.Outside != nil {
			for _, point := range testcase.Outside {
				if testcase.Poly.Contains(point) {
					t.Fatalf("Expected polygon %v to not contain point %v", testcase.Poly, point)
				}
			}
		}
	}
}
