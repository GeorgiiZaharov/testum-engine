PANDOC := pandoc

ifeq ($(MAKECMDGOALS),)
ALL_MD := $(wildcard *.md)
ALL_WIKI := $(ALL_MD:.md=.wiki)
endif

all: $(ALL_WIKI)

%.wiki: %.md
	@$(PANDOC) -f markdown -t mediawiki -o $@ $<

%.md:
	@$(MAKE) $(patsubst %.md,%.wiki,$@)

clean:
	@rm -f *.wiki

rebuild: clean all

list:
	@echo "Markdown файлы: $(wildcard *.md)"
	@echo "Wiki файлы: $(wildcard *.wiki)"

check:
	@$(PANDOC) --version || (echo "Pandoc не найден!" && exit 1)

.PHONY: all clean rebuild list check
