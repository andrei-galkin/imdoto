package yandexsearch

import "testing"

func TestGetImageItemFromJson(t *testing.T) {
	jsonImage := "{\"reqid\":\"1555697019666545-739689295744703832242977-man1-3752-IMG\",\"freshness\":\"normal\",\"preview\":[{\"url\":\"https://im0-tub-ru.yandex.net/i?id=1f8ed4bd9f436f087a34378eb18049ad-l&amp;n=13\",\"fileSizeInBytes\":1004020,\"w\":2560,\"h\":1600,\"origin\":{\"w\":2560,\"h\":1600,\"url\":\"https://img1.goodfon.ru/original/2560x1600/3/96/apple-blue-yabloko.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=340e966ab743310a034dc821cfef835c-l&amp;n=13\",\"fileSizeInBytes\":1164888,\"w\":2560,\"h\":1440,\"origin\":{\"w\":2560,\"h\":1440,\"url\":\"http://www.kartinkijane.ru/download.php?file=201304/2560x1440/kartinkijane.ru-25354.jpg\"},\"isMixedImage\":true}],\"dups\":[{\"url\":\"https://im0-tub-ru.yandex.net/i?id=948bb1c585c8459381b8b9419b9f1238-l&amp;n=13\",\"fileSizeInBytes\":1848219,\"w\":3000,\"h\":1875,\"origin\":{\"w\":3840,\"h\":2400,\"url\":\"https://www.jooomshaper.com/data/out/9/IMG_364652.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=89c77a75f68db90ec0136129ca3443c5-l&amp;n=13\",\"fileSizeInBytes\":1490824,\"w\":2880,\"h\":1800,\"origin\":{\"w\":2880,\"h\":1800,\"url\":\"https://www.wallpapersrc.com/img/3ff33dac4cfd46e88632a7d11fe84159/wood-textures-with-apple-logo-2880x1800.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=f0b8476c6240697c3a1430429eb1ab02-l&amp;n=13\",\"fileSizeInBytes\":991524,\"w\":2160,\"h\":1440,\"origin\":{\"w\":2160,\"h\":1440,\"url\":\"https://www.wallpapersrc.com/img/3ff33dac4cfd46e88632a7d11fe84159/wood-textures-with-apple-logo-2160x1440.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=72702ec8127748e09e82ade9c2e3bc68-l&amp;n=13\",\"fileSizeInBytes\":175390,\"w\":1024,\"h\":576,\"origin\":{\"w\":1024,\"h\":576,\"url\":\"https://www.desktopbackground.org/p/2012/04/16/375289_apple-wallpapers-hd-1920x1080-1370152_1920x1080_h.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=dee01d51f80094f996bd2e1ead50f94e-sr&amp;n=13\",\"fileSizeInBytes\":17811,\"w\":760,\"h\":428,\"origin\":{\"w\":380,\"h\":214,\"url\":\"https://www.artsfon.com/mini/201408/2943.jpg\"},\"isMixedImage\":true}],\"thumb\":{\"url\":\"//im0-tub-ru.yandex.net/i?id=7c4a17d4251cb1bc3c6c81af7e185f00&amp;n=13\",\"size\":{\"width\":480,\"height\":300},\"microImg\":\"\"},\"snippet\":{\"title\":\"Скачать обои Apple, Яблоко, Blue, раздел hi-tech в разрешени\",\"hasTitle\":true,\"text\":\"<b>Apple</b>, <b>Яблоко</b>, Blue. \",\"url\":\"https://www.goodfon.ru/download/apple-blue-yabloko/2560x1600/\",\"domain\":\"Goodfon.ru\",\"redirUrl\":\"http://yandex.ru/clck/jsredir?from=yandex.ru%3Bimages%2Fsearch%3Bimages%3B%3B&amp;text=&amp;etext=2127.BYq8vcENifRGimO0bAAhQPGg_Zy83dU9mfBj7NTm918.ac834183722a18635713d5570003231885e2cf38&amp;uuid=&amp;state=tid_Wvm4RM28ca_MiO4Ne9osTPtpHS9wicjEF5X7fRziVPIHCd9FyQ,,&amp;data=UlNrNmk5WktYejY4cHFySjRXSWhXRFUwd2xLN0F6SEx4NDdCZ3c5NlRqWlVGWWk5NUxJZnpKenBHNHZ1Ni1HQU42cWtJZ25GZDVuM2RUR2hvZE1SaG8ydGFvcUNjd2RhRGhBWGw1M0g0anpBQ0p5a3RFdTZaalFQMXJYblh4UVRqNXFFWlRwd0dZWlF3VklSNEVNOHFwenBNeEtBUVd0Tg,,&amp;sign=7ac288daafca48b607b2383cd921fabc&amp;keyno=0&amp;b64e=2&amp;l10n=ru\"},\"detail_url\":\"/images/search?pos=2&amp;img_url=https%3A%2F%2Fwww.wallpapers.net%2Fweb%2Fwallpapers%2Fgrey-apple-logo-wallpaper-for-desktop-mobiles%2F3840x2160.jpg&amp;text=apple&amp;rpt=simage\",\"img_href\":\"https://img1.goodfon.ru/original/2560x1600/3/96/apple-blue-yabloko.jpg\",\"useProxy\":false,\"pos\":2,\"id\":\"59eb43621eaab0b467c520eba07b3293\",\"rimId\":\"95f00eb5137e4fc41c93ebbebf53eb58\",\"docid\":\"ZB0C412756723E85F\",\"greenUrlCounterPath\":\"8.228.471.241.13.141\",\"counterPath\":\"thumb/normal\"}"
	expectedUrl := "www.jooomshaper.com/data/out/9/IMG_364652.jpg"
	img, err := GetImageItemFromJson(jsonImage)

	if err != nil {
		t.Errorf("Incorrect json")
	}

	if img.Dups[0].Origin.URL == expectedUrl {
		t.Errorf("Expected Image url is incorrect %s expected %s", img.Dups[0].Origin.URL, expectedUrl)
	}
}

func TestGetImageItemFromEmptyJson(t *testing.T) {
	jsonImage := " "
	img, err := GetImageItemFromJson(jsonImage)

	if err == nil && len(img.Dups[0].Origin.URL) == 0 {
		t.Errorf("GetImageItemFromJson fails")
	}
}

func TestGetFileFullName(t *testing.T) {
	folderPath := `c:`
	expectedFullName := `c:\59eb43621eaab0b467c520eba07b3293_IMG_364652.jpg`
	jsonImage := "{\"reqid\":\"1555697019666545-739689295744703832242977-man1-3752-IMG\",\"freshness\":\"normal\",\"preview\":[{\"url\":\"https://im0-tub-ru.yandex.net/i?id=1f8ed4bd9f436f087a34378eb18049ad-l&amp;n=13\",\"fileSizeInBytes\":1004020,\"w\":2560,\"h\":1600,\"origin\":{\"w\":2560,\"h\":1600,\"url\":\"https://img1.goodfon.ru/original/2560x1600/3/96/apple-blue-yabloko.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=340e966ab743310a034dc821cfef835c-l&amp;n=13\",\"fileSizeInBytes\":1164888,\"w\":2560,\"h\":1440,\"origin\":{\"w\":2560,\"h\":1440,\"url\":\"http://www.kartinkijane.ru/download.php?file=201304/2560x1440/kartinkijane.ru-25354.jpg\"},\"isMixedImage\":true}],\"dups\":[{\"url\":\"https://im0-tub-ru.yandex.net/i?id=948bb1c585c8459381b8b9419b9f1238-l&amp;n=13\",\"fileSizeInBytes\":1848219,\"w\":3000,\"h\":1875,\"origin\":{\"w\":3840,\"h\":2400,\"url\":\"https://www.jooomshaper.com/data/out/9/IMG_364652.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=89c77a75f68db90ec0136129ca3443c5-l&amp;n=13\",\"fileSizeInBytes\":1490824,\"w\":2880,\"h\":1800,\"origin\":{\"w\":2880,\"h\":1800,\"url\":\"https://www.wallpapersrc.com/img/3ff33dac4cfd46e88632a7d11fe84159/wood-textures-with-apple-logo-2880x1800.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=f0b8476c6240697c3a1430429eb1ab02-l&amp;n=13\",\"fileSizeInBytes\":991524,\"w\":2160,\"h\":1440,\"origin\":{\"w\":2160,\"h\":1440,\"url\":\"https://www.wallpapersrc.com/img/3ff33dac4cfd46e88632a7d11fe84159/wood-textures-with-apple-logo-2160x1440.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=72702ec8127748e09e82ade9c2e3bc68-l&amp;n=13\",\"fileSizeInBytes\":175390,\"w\":1024,\"h\":576,\"origin\":{\"w\":1024,\"h\":576,\"url\":\"https://www.desktopbackground.org/p/2012/04/16/375289_apple-wallpapers-hd-1920x1080-1370152_1920x1080_h.jpg\"},\"isMixedImage\":true},{\"url\":\"https://im0-tub-ru.yandex.net/i?id=dee01d51f80094f996bd2e1ead50f94e-sr&amp;n=13\",\"fileSizeInBytes\":17811,\"w\":760,\"h\":428,\"origin\":{\"w\":380,\"h\":214,\"url\":\"https://www.artsfon.com/mini/201408/2943.jpg\"},\"isMixedImage\":true}],\"thumb\":{\"url\":\"//im0-tub-ru.yandex.net/i?id=7c4a17d4251cb1bc3c6c81af7e185f00&amp;n=13\",\"size\":{\"width\":480,\"height\":300},\"microImg\":\"\"},\"snippet\":{\"title\":\"Скачать обои Apple, Яблоко, Blue, раздел hi-tech в разрешени\",\"hasTitle\":true,\"text\":\"<b>Apple</b>, <b>Яблоко</b>, Blue. \",\"url\":\"https://www.goodfon.ru/download/apple-blue-yabloko/2560x1600/\",\"domain\":\"Goodfon.ru\",\"redirUrl\":\"http://yandex.ru/clck/jsredir?from=yandex.ru%3Bimages%2Fsearch%3Bimages%3B%3B&amp;text=&amp;etext=2127.BYq8vcENifRGimO0bAAhQPGg_Zy83dU9mfBj7NTm918.ac834183722a18635713d5570003231885e2cf38&amp;uuid=&amp;state=tid_Wvm4RM28ca_MiO4Ne9osTPtpHS9wicjEF5X7fRziVPIHCd9FyQ,,&amp;data=UlNrNmk5WktYejY4cHFySjRXSWhXRFUwd2xLN0F6SEx4NDdCZ3c5NlRqWlVGWWk5NUxJZnpKenBHNHZ1Ni1HQU42cWtJZ25GZDVuM2RUR2hvZE1SaG8ydGFvcUNjd2RhRGhBWGw1M0g0anpBQ0p5a3RFdTZaalFQMXJYblh4UVRqNXFFWlRwd0dZWlF3VklSNEVNOHFwenBNeEtBUVd0Tg,,&amp;sign=7ac288daafca48b607b2383cd921fabc&amp;keyno=0&amp;b64e=2&amp;l10n=ru\"},\"detail_url\":\"/images/search?pos=2&amp;img_url=https%3A%2F%2Fwww.wallpapers.net%2Fweb%2Fwallpapers%2Fgrey-apple-logo-wallpaper-for-desktop-mobiles%2F3840x2160.jpg&amp;text=apple&amp;rpt=simage\",\"img_href\":\"https://img1.goodfon.ru/original/2560x1600/3/96/apple-blue-yabloko.jpg\",\"useProxy\":false,\"pos\":2,\"id\":\"59eb43621eaab0b467c520eba07b3293\",\"rimId\":\"95f00eb5137e4fc41c93ebbebf53eb58\",\"docid\":\"ZB0C412756723E85F\",\"greenUrlCounterPath\":\"8.228.471.241.13.141\",\"counterPath\":\"thumb/normal\"}"

	img, err := GetImageItemFromJson(jsonImage)
	fullName := GetFileFullName(img, folderPath)
	if err != nil {
		t.Errorf("Incorrect json")
	}

	if len(img.Dups[0].Origin.URL) == 0 {
		t.Errorf("Image url is empty")
	}

	if fullName != expectedFullName {
		t.Errorf("%s fullName is incorrect, expected %s", fullName, expectedFullName)
	}
}
