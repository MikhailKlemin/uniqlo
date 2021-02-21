package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pkg/errors"
)

func setMeasurements(page *rod.Page, height int, weight int) error {
	//#aria_uclw_headline
	var elem *rod.Element

	heightSel := "#uclw_form_height"
	weightSel := "#uclw_form_weight"
	cmSel := "#uclw_height_element > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1)"
	kgSel := "#uclw_weight_element > div:nth-child(1) > div:nth-child(1) > div:nth-child(2)"

	// setting cms
	err := rod.Try(func() {
		elem = page.MustSearch(cmSel)

	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot find input text at"))

	}
	elem.MustClick()

	// setting kgs
	err = rod.Try(func() {
		elem = page.MustSearch(kgSel)

	})
	if err != nil {

	}
	elem.MustClick()

	// setting height
	err = rod.Try(func() {
		elem = page.MustSearch(heightSel)

	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot find heightSel element at"))

	}

	elem.MustInput(fmt.Sprintf("%d", height))

	// setting weight
	err = rod.Try(func() {
		elem = page.MustSearch(weightSel)

	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot find weightSel text"))

	}
	elem.MustInput(fmt.Sprintf("%d", weight))

	page.WaitLoad()

	//click next
	page.MustSearch("#uclw_save_info_button").MustClick()
	page.WaitLoad()
	return nil

}

func setTummy(page *rod.Page, tummy int) error {
	err := rod.Try(func() {
		page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", tummy)).WaitStable(2 * time.Second)
		page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", tummy)).MustClick()
	})
	if err != nil {
		return errors.Wrap(err, "cannot set tummy")
	}
	page.WaitLoad()
	return nil
}

func setChest(page *rod.Page, chest int) error {
	err := rod.Try(func() {
		page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", chest)).WaitStable(2 * time.Second)
		page.MustSearch(fmt.Sprintf("#uclw_item_shape_%d", chest)).MustClick()
		page.WaitLoad()
	})

	if err != nil {
		return errors.Wrap(err, "cannot set chest")
	}
	return nil
}

func setAge(page *rod.Page, age int) error {
	err := rod.Try(func() {
		//enter age
		page.MustSearch(".uclw_input_text").MustInput(fmt.Sprintf("%d", age))
		//save button
		page.MustSearch("#uclw_save_info_button").MustClick()
		page.WaitLoad()
	})
	if err != nil {
		return errors.Wrap(err, "cannot find uclw_save_info_button during setAge ")
	}

	return nil
}

func setFit(page *rod.Page, fit int) error {
	fits := make(map[int]int)
	fits[-2] = 0
	fits[-1] = 50
	fits[0] = 100
	fits[1] = 150 //average
	fits[2] = 200
	fits[3] = 250
	fits[4] = 300

	/*
		var js = `myevent  = $.Event('mousedown');
		elem = $('.uclw_noUi-base');
		pos = elem.offset().left;
		myevent.clientX = pos;
		myevent.clientY = 0;
		elem.trigger(myevent);
		pos+10;`
	*/

	//fmt.Println(page.MustEval("10").Int())

	//fmt.Printf("offset %d \n", val)
	page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
	//page.MustEval(`myevent = $.Event('mousedown')`)

	left := page.MustEval("$('.uclw_noUi-base').offset().left").Int()
	left += fits[fit]
	top := page.MustEval("$('.uclw_noUi-base').offset().top").Int()
	mouse := page.Mouse

	mouse.Move(float64(left), float64(top), 2)
	mouse.MustDown(proto.InputMouseButtonLeft)
	mouse.MustUp(proto.InputMouseButtonLeft)
	time.Sleep(15 * time.Second)

	err := rod.Try(func() {
		fmt.Println("[INFO] Setting Fit\t", fit)
		page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
		//elem := page.MustSearch(".uclw_noUi-base")
		/*val := page.MustEval("pos = $('.uclw_noUi-base').offset().left").Int()
		fmt.Printf("val:%d\n", val)
		val = page.MustEval("pos").Int()
		fmt.Printf("val:%d\n", val)
		*/

		/*_, err := page.Eval(fmt.Sprintf(js))
		if err != nil {
			log.Fatal(err)
		}
		page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
		*/
		page.MustSearch("#uclw_save_info_button").MustClick()
		page.WaitLoad()
	})

	if err != nil {
		return errors.Wrap(err, "cannot find uclw_save_info_button during setFit")
	}

	return nil
}

func setFitOld(page *rod.Page) error {
	err := rod.Try(func() {
		page.MustSearch(".uclw_noUi-base").WaitStable(2 * time.Second)
		page.MustSearch("#uclw_save_info_button").MustClick()
		page.WaitLoad()
	})

	if err != nil {
		return errors.Wrap(err, "cannot find uclw_save_info_button during setFit")
	}

	return nil
}

func setLegsSize(page *rod.Page) error {
	err := rod.Try(func() {
		page.MustSearch(fmt.Sprintf("#uclw_item_tableletter_2")).WaitStable(2 * time.Second)
		page.MustSearch("#uclw_item_tableletter_2").MustClick()
		page.WaitLoad()
	})

	if err != nil {
		return errors.Wrap(err, "cannot find uclw_item_tableletter_2 during setLegs")
	}

	err = rod.Try(func() {
		//#uclw_item_lengthInput_skip
		page.MustSearch(fmt.Sprintf("#uclw_item_lengthInput_skip")).WaitStable(2 * time.Second)
		page.MustSearch("#uclw_item_lengthInput_skip").MustClick()
		page.WaitLoad()
	})

	if err != nil {
		return errors.Wrap(err, "cannot find uclw_item_lengthInput_skip during setLegs")
	}

	return nil

}

func getHeader(page *rod.Page) (header string, err error) {
	// #aria_uclw_headline
	//#aria_uclw_headline
	var elem *rod.Element

	err = rod.Try(func() {
		page.MustElement(`#aria_uclw_headline`).WaitStable(3 * time.Second)
	})

	if err != nil {
		return "", errors.Wrap(err, "Cannot find header on stable")
	}

	err = rod.Try(func() {
		elem = page.MustElement(`#aria_uclw_headline`)
	})
	if err != nil {
		return "", errors.Wrap(err, "Cannot find header")
	}

	header, err = elem.Text()
	if err != nil {
		return "", errors.Wrap(err, "Cannot get text from header")
	}

	//fmt.Println("INFO\t", strings.TrimSpace(header))
	return strings.TrimSpace(header), nil

}
