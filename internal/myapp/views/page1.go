package views

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/peterhellberg/hn"
	log "github.com/sirupsen/logrus"
	"hnclient/internal/myapp/useCases"
	"hnclient/pkg/gfxHelpers"
	"hnclient/pkg/widgets"
	"image/color"
	"math"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Page1 struct {
	theme material.Theme
	env   Enver

	topbar *TopBar
	title  string

	mainLayout      *layout.Stack
	itemsListLayout *layout.List

	storiesItems   map[int]*page1ListItem
	storiesLock    sync.RWMutex
	storiesUpdated chan struct{}
}

type Enver interface {
	GetInsets() layout.Inset
}

type page1ListItem struct {
	hn.Item

	site string
}

func NewPage1(th material.Theme, e Enver, title string) *Page1 {
	p := &Page1{
		theme:  th,
		env:    e,
		topbar: NewTopBar(th),
		title:  title,

		mainLayout: &layout.Stack{},
		itemsListLayout: &layout.List{
			Axis: layout.Vertical,
		},
		storiesUpdated: make(chan struct{}, 1),
	}

	go func() {
		log.Debug("requesting stories")
		in := make(chan useCases.IndexItem)

		p.storiesLock.Lock()
		p.storiesItems = make(map[int]*page1ListItem)
		p.storiesLock.Unlock()

		err := useCases.GetTopStories(50, in, time.Duration(15*time.Second))
		if err != nil {
			log.Error(err)
			return
		}

		for i := range in {
			p.storiesLock.Lock()
			li := page1ListItem{
				Item: *i.Item,
			}
			p.storiesItems[i.Index] = &li
			p.storiesLock.Unlock()
		}

		log.Debug("stories stored")
		p.storiesUpdated <- struct{}{}
	}()

	return p
}

func (p *Page1) Layout(gtx *layout.Context) {
	envInsets := p.env.GetInsets()

	baseInset := layout.Inset{
		Top:    unit.Max(gtx, unit.Dp(4), envInsets.Top),
		Right:  unit.Max(gtx, unit.Dp(12), envInsets.Right),
		Bottom: unit.Max(gtx, unit.Dp(4), envInsets.Bottom),
		Left:   unit.Max(gtx, unit.Dp(12), envInsets.Left),
	}
	baseRowInset := layout.Inset{
		Right: baseInset.Right,
		Left:  baseInset.Left,
	}
	statusMsgInset := layout.Inset{
		Top:    unit.Dp(12),
		Right:  baseRowInset.Right,
		Bottom: unit.Dp(12),
		Left:   baseRowInset.Left,
	}

	f := layout.Flex{Axis: layout.Vertical}

	topBarLayout := layout.Rigid(func() {
		p.topbar.Layout(gtx, p.env.GetInsets(), func() {
			lbl := material.H6(&p.theme, p.title)
			lbl.Color = gfxHelpers.RGB(0xffffff)
			lbl.Layout(gtx)
		})
	})

	// todo: how to get current label font from theme?
	//ls := p.theme.Shaper.LayoutString(text.Font{Variant: "Mono"}, 28, gtx.Constraints.Width.Max, "000")
	//var idxWidth int
	//if len(ls) > 0 {
	//	idxWidth = ls[0].Width.Floor()
	//}

	idxWidth := int(math.Floor(float64(unit.Dp(14 * 5).V)))

	messagesLayout := layout.Flexed(10, func() {
		p.storiesLock.RLock()
		defer p.storiesLock.RUnlock()

		if len(p.storiesItems) == 0 {
			in := statusMsgInset

			txt := "Loading..."

			in.Layout(gtx, func() {
				ll := material.Label(
					&p.theme,
					unit.Dp(12),
					txt,
				)
				ll.Color = gfxHelpers.RGB(0x666666)
				ll.Font.Style = text.Italic
				ll.Alignment = text.Middle
				ll.Layout(gtx)
			})
			return
		}

		p.itemsListLayout.Layout(gtx, len(p.storiesItems), func(index int) {
			item := p.storiesItems[index]

			// todo: improve item layout

			layout.Inset{
				Bottom: unit.Dp(1),
			}.Layout(gtx, func() {
				(&widgets.Background{
					Color: color.RGBA{
						R: 255,
						G: 102,
						B: 0,
						A: 0xaa,
					},
				}).Layout(gtx, func() {

					f := layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}

					var col1Height int

					col1 := layout.Rigid(func() {
						layout.Inset{
							Right:  baseRowInset.Left,
							Left:   baseRowInset.Left,
							Top:    unit.Dp(8),
							Bottom: unit.Dp(8),
						}.Layout(gtx, func() {
							// todo: width must be recalculated to 3 digits and set to all
							gtx.Constraints.Width.Min = idxWidth

							f := layout.Flex{
								Axis: layout.Vertical,
							}

							idx := layout.Rigid(func() {
								l := material.Label(
									&p.theme,
									unit.Dp(14),
									strconv.FormatInt(int64(index+1), 10),
								)
								l.Font.Variant = text.Variant("Mono")
								l.Alignment = text.Middle
								l.Layout(gtx)
							})

							score := layout.Rigid(func() {
								layout.Inset{
									Top:    unit.Dp(8),
									Bottom: unit.Dp(4),
								}.Layout(gtx, func() {

									if item == nil {
										return
									}

									var scoreStr string
									if item.Score > 999 {
										scoreStr = ">999"
									} else {
										scoreStr = strconv.FormatInt(int64(item.Score), 10)
									}

									l := material.Label(
										&p.theme,
										unit.Dp(10),
										scoreStr+"p",
									)
									l.Alignment = text.Middle
									l.Font.Variant = text.Variant("Mono")
									l.Layout(gtx)
								})
							})

							f.Layout(gtx, idx, score)

							col1Height = gtx.Dimensions.Size.Y
						})
					})

					col2 := layout.Rigid(func() {
						(&widgets.Background{
							Color: color.RGBA{
								R: 0xee,
								G: 0xee,
								B: 0xee,
								A: 0xff,
							},
							Inset: layout.Inset{
								Top:    unit.Dp(8),
								Right:  baseRowInset.Left,
								Bottom: unit.Dp(8),
								Left:   baseRowInset.Right,
							},
						}).Layout(gtx, func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max

							// make height min same as col1 rendered height
							if col1Height > gtx.Constraints.Height.Min {
								gtx.Constraints.Height.Min = col1Height
							}

							if item == nil {
								l := material.Label(&p.theme, unit.Dp(12), "...")
								l.Layout(gtx)
								return
							}

							col2rows := layout.Flex{Axis: layout.Vertical}

							col2row1 := layout.Rigid(func() {
								layout.Inset{
									Bottom: unit.Dp(8),
								}.Layout(gtx, func() {
									l := material.Label(&p.theme, unit.Dp(14), item.Title)
									l.Layout(gtx)
								})
							})
							col2row2 := layout.Rigid(func() {
								layout.Inset{
									Bottom: unit.Dp(4),
								}.Layout(gtx, func() {
									l := material.Label(&p.theme, unit.Dp(12), item.GetSite())
									l.Layout(gtx)
								})
							})
							col2row3 := layout.Rigid(func() {
								l := material.Label(&p.theme, unit.Dp(12),
									fmt.Sprintf("%s - %s",
										item.Time().Format("2006-01-02 15:04:05"),
										item.By,
									))
								l.Layout(gtx)
							})
							col2row4 := layout.Rigid(func() {
								var txt string
								if item.Deleted {
									txt = "Deleted"
								} else if item.Dead {
									txt = "Dead"
								}
								if txt != "" {
									layout.Inset{
										Top: unit.Dp(4),
									}.Layout(gtx, func() {
										l := material.Label(
											&p.theme,
											unit.Dp(12),
											txt)
										l.Color = color.RGBA{
											R: 0xff,
											G: 0x20,
											B: 0x20,
											A: 0xff,
										}
										l.Layout(gtx)
									})
								}
							})

							col2rows.Layout(gtx, col2row1, col2row2, col2row3, col2row4)
						})
					})

					f.Layout(gtx, col1, col2)
				})
			})
		})
	})

	f.Layout(gtx,
		topBarLayout,
		messagesLayout,
	)
}

//func (p *Page1) Event(gtx *layout.Context) (ev interface{}, redraw bool) {
//	select {
//	case <-p.storiesUpdated:
//		log.Debug("emitting redraw")
//		return nil, true
//
//
//	default:
//	}
//
//	return nil, false
//}

func (p *Page1) MustRedraw() <-chan struct{} {
	return p.storiesUpdated
}

func (i *page1ListItem) GetSite() string {
	if i.site != "" {
		return i.site
	}

	u, err := url.Parse(i.URL)
	if err != nil {
		return ""
	}

	i.site = u.Host

	return i.site
}
